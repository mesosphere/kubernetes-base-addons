import json
import subprocess
import os
import re
import pprint
import requests
import uuid
import ruamel.yaml

from functools import cmp_to_key


def compare_versions(a, b):
    pattern = r'\d+(\.\d+)*'
    m = re.search(pattern, a)
    if not m:
        return -1
    a = m.group(0)

    m = re.search(pattern, b)
    if not m:
        return 1
    b = m.group(0)

    a = a.split('.')
    b = b.split('.')

    for i in range(min(len(a), len(b))):
        if int(a[i]) > int(b[i]):
            return 1
        elif int(a[i]) < int(b[i]):
            return -1

    if len(a) > len(b):
        return -1
    elif len(a) < len(b):
        return 1
    return 0


def compare_yaml_files(a, b):
    pattern = r'([^-]+)(-)(\d+)'
    m = re.search(pattern, a)
    a_num = m.group(3)

    m = re.search(pattern, b)
    b_num = m.group(3)

    if int(a_num) > int(b_num):
        return 1
    elif int(a_num) < int(b_num):
        return -1
    return 0


def convert_repo_url(chart_repo_url, code_to_url, url_to_code):
    '''
    This function converts a repo url to a string that can be assigned as a repo name, which is required
    with the command "helm repo add".
    '''
    repo_code = chart_repo_url.split('://')[1]
    repo_code = repo_code.replace('/', '-')
    code_to_url[repo_code] = chart_repo_url
    url_to_code[chart_repo_url] = repo_code
    return repo_code


def update_annotations(loaded_yaml, info, app_version):
    res = compare_versions(info['app_version'], app_version)
    update_app_version = res != 0

    for k, v in loaded_yaml['metadata']['annotations'].items():
            if k.startswith('appversion') and update_app_version:
                loaded_yaml['metadata']['annotations'][k] = app_version
            elif k.startswith('catalog'):
                if update_app_version:
                    loaded_yaml['metadata']['annotations'][k] = app_version + '-1'
                else:
                    revision_number = 1
                    split_revision = v.split('-')
                    if len(split_revision) != 0:
                        try:
                            revision_number = int(split_revision[-1]) + 1
                        except ValueError:
                            pass
                    loaded_yaml['metadata']['annotations'][k] = '{}-{}'.format('-'.join(split_revision[:-1]), revision_number)
            elif k.startswith('docs'):
                minor_current_app_version = '.'.join(info['app_version'].split('.')[:2])
                minor_new_app_version = '.'.join(app_version.split('.')[:2])
                loaded_yaml['metadata']['annotations'][k] = v.replace(minor_current_app_version, minor_new_app_version)
            elif k.startswith('values'):
                m = re.search('.com/([^/]+/[^/]+)/([^/]+)', v)
                org_and_repo = m.group(1)
                old_sha = m.group(2)

                url = 'https://api.github.com/repos/{}/commits/master'.format(org_and_repo)
                headers = {'Authorization': 'token ' + os.environ['GITHUB_TOKEN']}
                r = requests.get(url, headers=headers)
                print(r.content)
                r.raise_for_status()
                new_sha = r.json()['sha'][:7]
                loaded_yaml['metadata']['annotations'][k] = v.replace(old_sha, new_sha)


def update_chart(chart, info):
    print('searching for {}/{}'.format(info['repo'], chart))
    search_cmd = ['helm', 'search', 'repo', '{}/{}'.format(info['repo'], chart)]
    print('Running search command: ' + str(search_cmd))
    sub = subprocess.run(search_cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    output = str(sub.stdout)
    column_titles = output.split('\\n')[0].split('\\t')

    # removing 2 first characters from column_titles[0] as they represent start of byte string (b')
    assert column_titles[0].strip()[2:] == 'NAME' and column_titles[1].strip() == 'CHART VERSION' and column_titles[2].strip() == 'APP VERSION'
    print('result from helm search: ' + output)
    try:
        result = output.split('\\n')[1].split('\\t')
        chart_version_from_search = result[1].strip()
        app_version = result[2].strip()
    except:
        raise Exception("Can't fetch latest version of chart {}. Output from helm search: {}".format(chart, sub.stdout.decode('utf-8')))

    res = compare_versions(info['chart_version'], chart_version_from_search)
    if res == -1:
        print('Newer chart version found: ' + chart_version_from_search + '\n')
        subprocess.run(['git', 'checkout', 'master'], check=True)
        random_string = uuid.uuid4().hex[:8]
        new_branch_name = 'bump-{}-{}'.format(chart, random_string)
        subprocess.run(['git', 'checkout', '-b', new_branch_name], check=True)

        with open(info['file_path'], 'r+') as stream:
            yaml = ruamel.yaml.YAML()
            yaml.preserve_quotes = True
            loaded = yaml.load(stream)
            loaded['spec']['chartReference']['version'] = chart_version_from_search
            if app_version != 'latest':
                update_annotations(loaded, info, app_version)
            stream.truncate(0)
            stream.seek(0)
            yaml.indent(sequence=4, offset=2)
            yaml.dump(loaded, stream)
            subprocess.run(["git", "commit", "-am", '"Bump {} to {}"'.format(chart, chart_version_from_search)], check=True)

        subprocess.run(['git', 'push', '-u', 'origin', new_branch_name], check=True)
        url = 'https://api.github.com/repos/mesosphere/kubernetes-base-addons/pulls'
        headers = {'Authorization': 'token ' + os.environ['GITHUB_TOKEN']}
        data = {
            'title': 'Automated chart bump {}-{}'.format(chart, chart_version_from_search),
            'head': new_branch_name,
            'base': 'master'
        }
        data = json.dumps(data).encode('utf-8')
        r = requests.post(url, data=data, headers=headers)
        print(r.content)
        r.raise_for_status()
    else:
        print('Chart version is already at the latest.\n')


def get_addon_dir():
    # Separating this in a function to enable mocking in unit tests
    return os.path.join(os.path.abspath(os.path.dirname(__file__)), '../../addons')

def main():
    # make sure github token is available
    os.environ['GITHUB_TOKEN']
    addons = {}
    code_to_url = {}
    url_to_code = {}
    addon_dir = get_addon_dir()
    for folder in os.listdir(addon_dir):
        subfolders = os.listdir(os.path.join(addon_dir, folder))
        subfolders.sort(key=cmp_to_key(compare_versions))
        latest_subfolder = ''
        if subfolders:
            latest_subfolder = subfolders[-1]
        yaml_files = [file for file in os.listdir(os.path.join(addon_dir, folder, latest_subfolder))]
        yaml_files.sort(key=cmp_to_key(compare_yaml_files))
        latest_yaml_file = yaml_files[-1]
        file_path = os.path.join(addon_dir, folder, latest_subfolder, latest_yaml_file)
        with open(file_path, 'r') as stream:
            yaml = ruamel.yaml.YAML()
            loaded = yaml.load(stream)

        chart_version = loaded['spec']['chartReference']['version']
        chart_name = loaded['spec']['chartReference']['chart']
        if 'stable' in chart_name:
            chart_name = chart_name.split('/')[1]
        chart_repo_url = loaded['spec']['chartReference'].get('repo', 'https://kubernetes-charts.storage.googleapis.com')

        app_version = ''
        for k, v in loaded['metadata']['annotations'].items():
            if k.startswith('appversion'):
                app_version = v

        repo_code = url_to_code.get(chart_repo_url)
        if repo_code is None:
            repo_code = convert_repo_url(chart_repo_url, code_to_url, url_to_code)

        addons[chart_name] = {
            'repo': repo_code,
            'chart_version': chart_version,
            'app_version': app_version,
            'file_path': file_path,
        }

    pprint.pprint(addons, indent=8)

    for code, url in code_to_url.items():
        subprocess.run(['helm', 'repo', 'add', code, url], check=True)

    subprocess.run(['helm', 'repo', 'list'], check=True)
    subprocess.run(['helm', 'repo', 'update'], check=True)

    for chart, info in addons.items():
        update_chart(chart, info)


if __name__ == '__main__':
    main()
