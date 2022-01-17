# Copyright (C) 2022 Hatching B.V.
# All rights reserved.

import pytest

from mock import patch, Mock
from requests import Session

import triage

class TestPaginator:

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    def test_no_data(self, s, r):
        c = triage.Client("token")
        m = Mock()
        m.json = Mock(return_value={
            "data": None,
            "next": None
        })
        s.return_value = m

        p = c.profiles()
        with pytest.raises(StopIteration):
            next(p)

        r.assert_called_with(
            "GET",
            "/v0/profiles?limit=20",
            None
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    def test_next(self, s, r):
        c = triage.Client("token")
        m = Mock()
        m.json = Mock(return_value={
            "data": [{"name": "test"}],
            "next": 1
        })
        s.return_value = m

        p = c.profiles()
        assert next(p)["name"] == "test"
        with pytest.raises(StopIteration):
            next(p)

        r.assert_called_with(
            "GET",
            "/v0/profiles?limit=20&offset=1",
            None
        )

    @patch('triage.Client._new_request')
    @patch.object(Session, 'send')
    def test_no_next(self, s, r):
        c = triage.Client("token")
        m = Mock()
        m.json = Mock(return_value={
            "data": [{"name": "test"}],
            "next": None
        })
        s.return_value = m

        p = c.profiles()
        next(p)
        with pytest.raises(StopIteration):
            next(p)

        r.assert_called_with(
            "GET",
            "/v0/profiles?limit=20",
            None
        )