#!/usr/bin/env python3
"""Convert viper keys to this repository's env naming convention."""

import argparse


def to_env_name(viper_key: str) -> str:
    return viper_key.strip().upper().replace(".", "__")


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Convert a viper key like host.api_port to HOST__API_PORT"
    )
    parser.add_argument("viper_key", help="viper key path (example: host.api_port)")
    args = parser.parse_args()
    print(to_env_name(args.viper_key))


if __name__ == "__main__":
    main()
