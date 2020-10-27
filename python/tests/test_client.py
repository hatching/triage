# Copyright (C) 2020 Hatching B.V.
# All rights reserved.

import triage
import io
import pytest

from mock import patch, Mock

def helper(method_, path_, data_, mock_urlopen, return_value):
    def _new_request(method, path, data=None):
        assert method == method_
        assert path == path_
        if data_:
            assert data.read() == data_
        else:
            assert data is None
        return Mock()

    c = triage.Client("token")
    c._new_request = _new_request
    mock_urlopen.return_value.__enter__.return_value.read.\
        return_value = return_value
    return c

class TestSampleAction:
    @pytest.mark.skip(reason="Currently not checking form data")
    @patch('triage.client.urlopen')
    def test_submit_file(self, mock_urlopen):
        c = helper(
            "POST", "/v0/samples", None,
            mock_urlopen, '{}'
        )
        c.submit_sample_file("test", io.StringIO("file"))

    @patch('triage.client.urlopen')
    def test_submit_url(self, mock_urlopen):
        c = helper(
            "POST", "/v0/samples",
                '{"kind": "url", "url": "http://9gag.com"'
                ', "interactive": false, "profiles": []}',
            mock_urlopen, '{}'
        )
        c.submit_sample_url("http://9gag.com")

    @patch('triage.client.urlopen')
    def test_profile(self, mock_urlopen):
        c = helper(
            "POST", "/v0/samples/sample1/profile",
                '{"auto": false, "profiles": []}',
            mock_urlopen, '{}'
        )
        c.set_sample_profile("sample1", [])

    @patch('triage.client.urlopen')
    def test_profile_automatic(self, mock_urlopen):
        c = helper(
            "POST", "/v0/samples/sample1/profile",
                '{"auto": true, "pick": []}',
            mock_urlopen, '{}'
        )
        c.set_sample_profile_automatically("sample1")

class TestReport:
    @patch('triage.client.urlopen')
    def test_owned_samples(self, mock_urlopen):
        c = helper(
            "GET", "/v0/samples?subset=owned&limit=20", None,
            mock_urlopen, '{"data": [{"name": "test"}]}'
        )
        for i in c.owned_samples():
            assert i["name"] == "test"

    @patch('triage.client.urlopen')
    def test_public_samples(self, mock_urlopen):
        c = helper(
            "GET", "/v0/samples?subset=public&limit=20", None,
            mock_urlopen, '{"data": [{"name": "test"}]}'
        )
        for i in c.public_samples():
            assert i["name"] == "test"

    @patch('triage.client.urlopen')
    def test_search(self, mock_urlopen):
        c = helper(
            "GET", "/v0/search?query=NOT+family%3Aemotet&limit=20", None,
            mock_urlopen, '{"data": [{"name": "test"}]}'
        )
        for i in c.search("NOT family:emotet"):
            assert i["name"] == "test"

        c = helper(
            "GET", "/v0/search?query=NOT+family%3Aemotet&limit=200", None,
            mock_urlopen, '{"data": [{"name": "test"}]}'
        )
        for i in c.search("NOT family:emotet", 1000):
            assert i["name"] == "test"

    @patch('triage.client.urlopen')
    def test_sample(self, mock_urlopen):
        c = helper(
            "GET", "/v0/samples/sample1", None,
            mock_urlopen, '{}'
        )
        c.sample_by_id("sample1")

    @patch('triage.client.urlopen')
    def test_delete(self, mock_urlopen):
        c = helper(
            "DELETE", "/v0/samples/sample1", None,
            mock_urlopen, '{}'
        )
        c.delete_sample("sample1")

    @patch('triage.client.urlopen')
    def test_static(self, mock_urlopen):
        c = helper(
            "GET", "/v0/samples/sample1/reports/static", None,
            mock_urlopen, '{}'
        )
        c.static_report("sample1")

    @patch('triage.client.urlopen')
    def test_overview(self, mock_urlopen):
        c = helper(
            "GET", "/v0/samples/sample1/overview.json", None,
            mock_urlopen, '{}'
        )
        c.overview_report("sample1")

    @patch('triage.client.urlopen')
    def test_task(self, mock_urlopen):
        c = helper(
            "GET", "/v0/samples/sample1/task1/report_triage.json", None,
            mock_urlopen, '{}'
        )
        c.task_report("sample1", "task1")

class TestFile:
    @patch('triage.client.urlopen')
    def test_file(self, mock_urlopen):
        c = helper(
            "GET", "/v0/samples/sample1/task1/file1", None,
            mock_urlopen, '{}'
        )
        c.sample_task_file("sample1", "task1", "file1")

    @patch('triage.client.urlopen')
    def test_tar(self, mock_urlopen):
        c = helper(
            "GET", "/v0/samples/test1/archive", None,
            mock_urlopen, '{}'
        )
        c.sample_archive_tar("test1")

    @patch('triage.client.urlopen')
    def test_zip(self, mock_urlopen):
        c = helper(
            "GET", "/v0/samples/test1/archive.zip", None,
            mock_urlopen, '{}'
        )
        c.sample_archive_zip("test1")

class TestProfile:
    @patch('triage.client.urlopen')
    def test_list_profiles(self, mock_urlopen):
        c = helper(
            "GET", "/v0/profiles?limit=20", None,
            mock_urlopen, '{"data": [{"name": "test"}]}'
        )
        for x in c.profiles():
            assert x["name"] == "test"

    @patch('triage.client.urlopen')
    def test_delete_profile(self, mock_urlopen):
        c = helper(
            "DELETE", "/v0/profiles/meme", None,
            mock_urlopen, '{}'
        )
        c.delete_profile("meme")

    @patch('triage.client.urlopen')
    def test_create_profile(self, mock_urlopen):
        c = helper(
            "POST", "/v0/profiles",
                '{"name": "meme", "tags": ["yes1", "yes2"], '
                '"network": "drop", "timeout": 30}',
            mock_urlopen, '{}'
        )
        c.create_profile("meme", ["yes1","yes2"], "drop", 30)
