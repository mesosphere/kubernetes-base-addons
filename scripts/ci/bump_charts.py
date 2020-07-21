import subprocess
import os
import re
import yaml
import pprint
import random
import string
import threading
import collections

from functools import cmp_to_key


def compare_versions(a, b):
    a = a.split('.')
    b = b.split('.')

    for i in range(min(len(a), len(b))):
        if int(a[i]) > int(b[i]):
            return 1
        elif int(a[i]) < int(b[i]):
            return -1

    if len(a) > len(b):
        return 1
    elif len(a) < len(b):
        return -1
    return 0


def compare_subfolders(a, b):
    pattern = '\d+(\.\d+)*'
    m = re.search(pattern, a)
    short_a = m.group(0)

    m = re.search(pattern, b)
    short_b = m.group(0)

    return compare_versions(a, b)


def compare_yaml_files(a, b):
    pattern = '([^-]+)(-)(\d+)'
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
    repo_code = chart_repo_url.split('://')[1]
    repo_code = repo_code.replace('/', '-')
    code_to_url[repo_code] = chart_repo_url
    url_to_code[chart_repo_url] = repo_code
    return repo_code


def update_chart(chart, info):
    print('searching for {}/{}'.format(info['repo'], chart))
    search_cmd = ['helm', 'search', 'repo', '{}/{}'.format(info['repo'], chart)]
    print('Running search command: ' + str(search_cmd))
    sub = subprocess.run(search_cmd, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
    output = str(sub.stdout)
    result = output.split('\\n')[1]
    print('result from helm search: ' + output)
    try:
        chart_version = result.split('\\t')[1].strip()
    except:
        raise Exception("Can't fetch latest version of chart {}. Output from helm search: {}".format(chart, sub.stdout.decode('utf-8')))
    
    pattern = '\d+(\.\d+)*'
    m = re.search(pattern, info['version'])
    m1 = re.search(pattern, chart_version)
    res = compare_versions(m.group(0), m1.group(0))
    if res == -1:
        print('Newer version found: ' + chart_version + '\n')
        with open(info['file_path'], 'r+') as stream:
            original_content = stream.read()
            start = original_content.find('chartReference:')
            if start == -1:
                raise Exception("Error parsing chart yaml: chartReference not found")
            content = original_content[start:]
            index = content.find('version:')
            if index == -1:
                raise Exception("Error parsing chart yaml: chartReference version not found")
            start += index
            start += len('version:')
            content = original_content[start:]
            end = content.find('\n')
            if end == -1:
                raise Exception("Error parsing chart yaml: end of line not found")

            updated_content = original_content[:start] + ' ' + chart_version + content[end:]
            stream.seek(0)
            stream.write(updated_content)
    else:
        print('Version is already the latest.\n')


if __name__ == '__main__':
    addons = {}
    code_to_url = {}
    url_to_code = {}
    addon_dir = os.path.join(os.path.abspath(os.path.dirname(__file__)), '../../addons')
    for folder in os.listdir(addon_dir):
        subfolders = os.listdir(os.path.join(addon_dir, folder))
        subfolders.sort(key=cmp_to_key(compare_subfolders))
        latest_subfolder = subfolders[-1]
        yaml_files = [file for file in os.listdir(os.path.join(addon_dir, folder, latest_subfolder))]
        yaml_files.sort(key=cmp_to_key(compare_yaml_files))
        file_path = os.path.join(addon_dir, folder, latest_subfolder, yaml_files[-1])
        with open(file_path, 'r') as stream:
            loaded = yaml.load(stream, Loader=yaml.FullLoader)

        chart_version = loaded['spec']['chartReference']['version']
        chart_name = loaded['spec']['chartReference']['chart']
        if 'stable' in chart_name:
            chart_name = chart_name.split('/')[1]
        chart_repo_url = loaded['spec']['chartReference'].get('repo', 'https://kubernetes-charts.storage.googleapis.com')

        repo_code = url_to_code.get(chart_repo_url)
        if repo_code is None:
            repo_code = convert_repo_url(chart_repo_url, code_to_url, url_to_code)

        addons[chart_name] = {
            'repo': repo_code,
            'version': chart_version,
            'file_path': file_path
        }

    pprint.pprint(addons, indent=8)

    for code, url in code_to_url.items():
        subprocess.run(['helm', 'repo', 'add', code, url], check=True)

    subprocess.run(['helm', 'repo', 'list'], check=True)
    subprocess.run(['helm', 'repo', 'update'], check=True)

    for chart, info in addons.items():
        update_chart(chart, info)
