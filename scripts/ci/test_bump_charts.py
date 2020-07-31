import bump_charts

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
