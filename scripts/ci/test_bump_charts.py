import os
import bump_charts
import shutil
from unittest.mock import Mock

def test_compare_yaml_files():
    res = bump_charts.compare_yaml_files('ops-portal-9.yaml', 'ops-portal-10.yaml')
    assert res == -1

    res = bump_charts.compare_yaml_files('ops-portal-9.yaml', 'ops-portal-9.yaml')
    assert res == 0

    res = bump_charts.compare_yaml_files('ops-portal-10.yaml', 'ops-portal-9.yaml')
    assert res == 1


def test_compare_versions():
    res = bump_charts.compare_versions('1.1.1', '1.1.12')
    assert res == -1

    res = bump_charts.compare_versions('1.1.1', '1.1')
    assert res == -1

    res = bump_charts.compare_versions('', '1.1')
    assert res == -1

    res = bump_charts.compare_versions('1.1.1', '1.1.1')
    assert res == 0

    res = bump_charts.compare_versions('1.1.1', '1.2.1')
    assert res == -1


def test_convert_repo_url():
    res = bump_charts.convert_repo_url('https://kubernetes-charts.storage.googleapis.com', {}, {})
    assert res == 'kubernetes-charts.storage.googleapis.com'


def test_main(tmp_path):
    bump_charts.os.environ = {'GITHUB_TOKEN': ''}

    # mocking addon dir
    new_tmp_path = tmp_path / 'cert-manager'
    new_tmp_path.mkdir()
    new_tmp_path = new_tmp_path / '0.10.x'
    new_tmp_path.mkdir()
    dir_path = os.path.join(os.path.abspath(os.path.dirname(__file__)), 'unit_test_yaml_samples')
    file1_path = os.path.join(dir_path, 'test-certmanager-1.yaml')
    file2_path = os.path.join(dir_path, 'test-certmanager-12.yaml')
    shutil.copy(file1_path, new_tmp_path)
    shutil.copy(file2_path, new_tmp_path)
    bump_charts.get_addon_dir = Mock(return_value=tmp_path)

    '''mocking:
    r = requests.get(args*).json()['sha'][:7]
    print(r.content)
    r.raise_for_status()
    '''
    mock_response = Mock(return_value=Mock(json=Mock(return_value={'sha': '1a2b3c4'})))
    mock_requests_get = mock_response
    mock_response.raise_for_status = Mock()
    mock_response.content = ''
    bump_charts.requests.get = mock_requests_get

    # same but for requests.post
    mock_requests_post = Mock()
    mock_requests_post.raise_for_status = Mock()
    mock_requests_post.content = ''
    bump_charts.requests.post = mock_requests_post

    # mocking: subprocess.run().stdout
    subprocess_mock = Mock()
    subprocess_mock.stdout = "NAME \tCHART VERSION\tAPP VERSION\tDESCRIPTION \nmesosphere.github.io-charts-staging/cert-manager-setup\t0.1.14 \t0.10.1 \tInstall cert-manager and optionally add a ClusterIssuer\n".encode('utf-8')
    bump_charts.subprocess.run = Mock(return_value=subprocess_mock)

    new_file1_path = new_tmp_path / 'test-certmanager-1.yaml'

    with open(file1_path, 'r') as f:
        file1_content_before_run = f.read()

    bump_charts.main()

    with open(new_file1_path, 'r') as f:
        file1_content_after_run = f.read()

    # test-certmanager-1.yaml should remain unchanged because it is not the latest file for certmanager
    assert file1_content_before_run == file1_content_after_run

    with open(new_tmp_path / 'test-certmanager-12.yaml', 'r') as f:
        file2_content_after_run = f.read()

    with open(os.path.join(dir_path, 'test-updated-certmanager-12.yaml'), 'r') as f:
        file2_updated_reference_content = f.read()

    # test-updated-certmanager-12.yaml has updated annotations and chart version of test-certmanager-12.yaml
    assert file2_content_after_run == file2_updated_reference_content
