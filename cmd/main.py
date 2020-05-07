#!/usr/bin/env python3

import argparse
from pathlib import Path
import sys
from dotenv import load_dotenv

# do this so we can call files in another folder
sys.path.append('./')
import internal_py.setup as setup

def main():
    # load configuration
    env_path = Path('.') / 'configs' / '.env'
    load_dotenv(dotenv_path=env_path)

    # setup logger
    setup.setup_logger()

    # run appropriate app
    flags = select_app()

    print("LOL")
    print(flags)


def select_app():
    # Construct the argument parser
    parser = argparse.ArgumentParser()

    # Add the arguments to the parser
    parser.add_argument("-news", action='store_true', help="Run Newspaper3k server")
    parser.add_argument("-text", action='store_true', help="Run text processing server")
    flags = parser.parse_args()

    return flags


# startup server
if __name__ == "__main__":
    main()
    