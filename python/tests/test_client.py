# Copyright (C) 2020-2022 Hatching B.V.
# All rights reserved.

import io

from mock import patch, Mock
from requests import Session

import triage

class BytesArg(object):
    def __eq__(a, b):
        return isinstance(b, io.BytesIO)

class MultiPartArg(object):
    def __eq__(a, b):
        return "multipart/form-data" in b["Content-Type"]

class TestSampleAction:
    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_submit_file(self, m, s, r):
        c = triage.Client("token")
        c.submit_sample_file("test", io.StringIO("file"))
        r.assert_called_with(
            "POST",
            "/v0/samples",
            b=BytesArg(),
            headers=MultiPartArg()
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_submit_url(self, m, s, r):
        c = triage.Client("token")
        c.submit_sample_url("http://9gag.com")
        r.assert_called_with(
            "POST",
            "/v0/samples",
            {
                'kind': 'url',
                'url': 'http://9gag.com',
                'interactive': False,
                'profiles': []
            },
            headers={'Content-Type': 'application/json'}
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_profile(self, m, s, r):
        c = triage.Client("token")
        c.set_sample_profile("sample1", [])
        r.assert_called_with(
            "POST",
            "/v0/samples/sample1/profile",
            {
                'auto': False,
                'profiles': []
            },
            headers={'Content-Type': 'application/json'}
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_profile_automatic(self, m, s, r):
        c = triage.Client("token")
        c.set_sample_profile_automatically("sample1")
        r.assert_called_with(
            "POST",
            "/v0/samples/sample1/profile",
            {
                'auto': True,
                'pick': []
            },
            headers={'Content-Type': 'application/json'}
        )


class TestReport:
    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_org_samples(self, m, s, r):
        c = triage.Client("token")
        m = Mock()
        m.json = Mock(return_value={
            "data": [{"name": "test"}],
            "next": None
        })
        s.return_value = m

        for i in c.org_samples():
            assert i["name"] == "test"
        m.json.assert_called_once()

        r.assert_called_with(
            "GET",
            "/v0/samples?subset=org&limit=20",
            None
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_owned_samples(self, m, s, r):
        c = triage.Client("token")
        m = Mock()
        m.json = Mock(return_value={
            "data": [{"name": "test"}],
            "next": None
        })
        s.return_value = m

        for i in c.owned_samples():
            assert i["name"] == "test"
        m.json.assert_called_once()

        r.assert_called_with(
            "GET",
            "/v0/samples?subset=owned&limit=20",
            None
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_public_samples(self, m, s, r):
        c = triage.Client("token")
        m = Mock()
        m.json = Mock(return_value={
            "data": [{"name": "test"}],
            "next": None
        })
        s.return_value = m

        for i in c.public_samples():
            assert i["name"] == "test"
        m.json.assert_called_once()

        r.assert_called_with(
            "GET",
            "/v0/samples?subset=public&limit=20",
            None
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_search(self, m, s, r):
        c = triage.Client("token")
        m = Mock()
        m.json = Mock(return_value={
            "data": [{"name": "test"}],
            "next": None
        })
        s.return_value = m

        for i in c.search("NOT family:emotet"):
            assert i["name"] == "test"
        m.json.assert_called_once()

        r.assert_called_with(
            "GET",
            "/v0/search?query=NOT%20family%3Aemotet&limit=20",
            None
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_search_pagination(self, m, s, r):
        c = triage.Client("token")
        m = Mock()
        m.json = Mock(return_value={
            "data": [{"name": "test"}],
            "next": None
        })
        s.return_value = m

        for i in c.search("NOT family:emotet", 1000):
            assert i["name"] == "test"
        m.json.assert_called_once()

        r.assert_called_with(
            "GET",
            "/v0/search?query=NOT%20family%3Aemotet&limit=200",
            None
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_sample(self, m, s, r):
        c = triage.Client("token")
        c.sample_by_id("sample1")
        r.assert_called_with(
            "GET",
            "/v0/samples/sample1",
            None
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_delete(self, m, s, r):
        c = triage.Client("token")
        c.delete_sample("sample1")
        r.assert_called_with(
            "DELETE",
            "/v0/samples/sample1",
            None
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_static(self, m, s, r):
        c = triage.Client("token")
        c.static_report("sample1")
        r.assert_called_with(
            "GET",
            "/v0/samples/sample1/reports/static",
            None
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_overview(self, m, s, r):
        c = triage.Client("token")
        c.overview_report("sample1")
        r.assert_called_with(
            "GET",
            "/v1/samples/sample1/overview.json",
            None
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_task(self, m, s, r):
        c = triage.Client("token")
        c.task_report("sample1", "task1")
        r.assert_called_with(
            "GET",
            "/v0/samples/sample1/task1/report_triage.json",
            None
        )

class TestFile:
    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_file(self, m, s, r):
        c = triage.Client("token")
        c.sample_task_file("sample1", "task1", "file1")
        r.assert_called_with(
            "GET",
            "/v0/samples/sample1/task1/file1",
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_tar(self, m, s, r):
        c = triage.Client("token")
        c.sample_archive_tar("test1")
        r.assert_called_with(
            "GET",
            "/v0/samples/test1/archive",
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_zip(self, m, s, r):
        c = triage.Client("token")
        c.sample_archive_zip("test1")
        r.assert_called_with(
            "GET",
            "/v0/samples/test1/archive.zip",
        )

class TestProfile:
    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_list_profiles(self, m, s, r):
        c = triage.Client("token")
        m = Mock()
        m.json = Mock(return_value={
            "data": [{"name": "test"}],
            "next": None
        })
        s.return_value = m

        for x in c.profiles():
            assert x["name"] == "test"
        m.json.assert_called_once()

        r.assert_called_with(
            "GET",
            "/v0/profiles?limit=20",
            None
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_delete_profile(self, m, s, r):
        c = triage.Client("token")
        c.delete_profile("meme")
        r.assert_called_with(
            "DELETE",
            "/v0/profiles/meme",
            None
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    @patch.object(Session, 'merge_environment_settings')
    def test_create_profile(self, m, s, r):
        c = triage.Client("token")
        c.create_profile("meme", ["yes1","yes2"], "drop", 30)
        r.assert_called_with(
            "POST",
            "/v0/profiles",
            {
                "name": "meme",
                "tags": ["yes1", "yes2"],
                "network":
                "drop", "timeout": 30
            },
            headers={'Content-Type': 'application/json'}
        )
