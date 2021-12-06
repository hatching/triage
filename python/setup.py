# Copyright (C) 2020-2021 Hatching B.V.
# All rights reserved.

import setuptools

with open("README.md", "r") as fh:
    long_description = fh.read()

setuptools.setup(
    name="hatching-triage",
    version="0.1.5",
    author="Hatching B.V.",
    author_email="info@hatching.io",
    description="API client and CLI for Hatching Triage",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://hatching.io/",
    packages=setuptools.find_packages(),
    entry_points={
        "console_scripts": [
            "triage = cli.triage:cli",
        ],
    },
    install_requires=[
        "appdirs==1.4.4",
        "click==8.0.3",
        "requests==2.25.1"
    ],
    python_requires='>=3.6',
)
