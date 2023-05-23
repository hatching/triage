# Copyright (C) 2020-2023 Hatching B.V.
# All rights reserved.

import setuptools

with open("README.md", "r") as fh:
    long_description = fh.read()

# Execute this so that __version__ is set. Don't import it, this causes
# installation fails on missing requirements.
with open("triage/__version__.py", "rb") as fp:
    exec(fp.read())

setuptools.setup(
    name="hatching-triage",
    version=__version__,
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
        "click>=8.0.3, <8.2",
        "requests>=2.25.1, <3"
    ],
    python_requires='>=3.6',
)
