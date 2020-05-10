#!/usr/bin/env python3
"""Starts program and handles input"""

import argparse
from pathlib import Path
import sys
from dotenv import load_dotenv

# import like this so we can call files in another folder
sys.path.append('./')
import internal_py.setup as setup
import internal_py.newspaper as newspaper

def main():
    """Performs general setup and call appropriate application"""

    # load configuration
    env_path = Path('.') / 'configs' / '.env'
    load_dotenv(dotenv_path=env_path)

    # setup logger
    setup.setup_logger()

    # run appropriate app
    flags = select_app()

    print("LOL")
    print(flags)
    newspaper.run()



def select_app():
    """Parses input flags to determine which program to run"""

    # Construct the argument parser
    parser = argparse.ArgumentParser()

    # Add the arguments to the parser
    parser.add_argument("-news", action='store_true', help="Run Newspaper3k server")
    parser.add_argument("-text", action='store_true', help="Run text processing server")
    flags = parser.parse_args()

    return flags

# Default program entry point
if __name__ == "__main__":
    main()
    