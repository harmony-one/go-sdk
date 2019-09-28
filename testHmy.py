#!/usr/bin/env python
import subprocess
import pexpect
import os
import shutil
import json
import logging
import inspect
import datetime
import sys
import random

# To flag a failure, make sure to exit with a non-zero exit code.

ENVIRONMENT = {}
ADDRESSES = {}
KEYSTORE_PATH = ""
KEYS_ADDED = set()


def log(message, error=True):
    func = inspect.currentframe().f_back.f_code
    final_msg = "(%s:%i) %s" % (
        func.co_name,
        func.co_firstlineno,
        message
    )
    if error:
        logging.warning(final_msg)
        print(f"[ERROR] {final_msg}")
    else:
        print(final_msg)


def load_environment():
    global ENVIRONMENT
    try:
        # Requires the updated 'setup_bls_build_flags.sh'
        env_raw = subprocess.check_output("source ../harmony/scripts/setup_bls_build_flags.sh -v", shell=True)
        ENVIRONMENT = json.loads(env_raw)
        ENVIRONMENT["HOME"] = os.environ.get("HOME")
    except json.decoder.JSONDecodeError as _:
        log(f"[Critical] Could not parse environment variables from setup_bls_build_flags.sh")
        sys.exit(-1)


def load_addresses():
    global ADDRESSES
    try:
        response = subprocess.check_output(["hmy", "keys", "list"], env=ENVIRONMENT).decode()
    except subprocess.CalledProcessError as err:
        log(f"[Critical] Could not list keys.\n"
            f"Got exit code {err.returncode}. Msg: {err.output}")
        sys.exit(-1)

    lines = response.split("\n")
    if "NAME" not in lines[0] or "ADDRESS" not in lines[0]:
        log(f"[Critical] Name or Address not found on first line if key list.")
        sys.exit(-1)

    for line in lines[1:]:
        if not line:
            continue
        try:
            name, address = line.split("\t")
        except ValueError:
            log(f"[Critical] Unexpected key list format.")
            sys.exit(-1)
        ADDRESSES[name.strip()] = address


def load_keystore_directory():
    global KEYSTORE_PATH
    try:
        response = subprocess.check_output(["hmy", "keys", "location"], env=ENVIRONMENT).decode().strip()
    except subprocess.CalledProcessError as err:
        log(f"[Critical] Could not get keystore path.\n"
            f"Got exit code {err.returncode}. Msg: {err.output}")
        sys.exit(-1)
    if not os.path.exists(response):
        log(f"[Critical] '{response}' is not a valid path")
        sys.exit(-1)
    KEYSTORE_PATH = response


def delete_from_keystore_by_name(name):
    log(f"[KEY DELETE] Removing {name} from keystore at {KEYSTORE_PATH}", error=False)
    key_file_path = f"{KEYSTORE_PATH}/{name}"

    try:
        shutil.rmtree(key_file_path)
    except shutil.Error as e:
        log(f"[KEY DELETE] Failed to delete dir: {key_file_path}\n"
            f"Exception: {e}")
        return

    del ADDRESSES[name]


def get_address_from_name(name):
    if name in ADDRESSES:
        return ADDRESSES[name]
    else:
        load_addresses()
        return ADDRESSES.get(name, None)


def test_mnemonics_from_sdk():
    log("Testing...", error=False)
    with open('testHmyConfigs/sdkMnemonics.json') as f:
        sdk_mnemonics = json.load(f)
        if not sdk_mnemonics:
            log("Could not load reference data.")
            return False

    passed = True
    for test in sdk_mnemonics["data"]:
        index = test["index"]
        if index != 0:  # CLI currently uses a hardcoded index of 0.
            continue

        mnemonic = test["phrase"]
        correct_address = test["addr"]
        address_name = f'testHmyAcc_{random.randint(0,1e9)}'
        while address_name in ADDRESSES:
            address_name = f'testHmyAcc_{random.randint(0,1e9)}'

        try:
            hmy = pexpect.spawn('./hmy', ['keys', 'add', address_name, '--recover', '--passphrase'], env=ENVIRONMENT)
            hmy.expect("Enter passphrase for account\r\n")
            hmy.sendline("")
            hmy.expect("Repeat the passphrase:\r\n")
            hmy.sendline("")
            hmy.expect("Enter mnemonic to recover keys from\r\n")
            hmy.sendline(mnemonic)
            hmy.wait()
            hmy.expect(pexpect.EOF)
        except pexpect.ExceptionPexpect as e:
            log(f"Exception occurred when adding a key with mnemonic."
                f"\nException: {e}")
            passed = False

        hmy_address = get_address_from_name(address_name)
        if hmy_address != correct_address or hmy_address is None:
            log(f"address does not match sdk's address. \n"
                f"\tMnemonic: {mnemonic}\n"
                f"\tCorrect address: {correct_address}\n"
                f"\tCLI address: {hmy_address}")
            passed = False
        else:
            KEYS_ADDED.add(address_name)

    return passed


if __name__ == "__main__":
    logging.basicConfig(filename="testHmy.log", filemode='a', format="%(message)s")
    logging.warning(f"[{datetime.datetime.now()}] {'=' * 10}")
    load_environment()
    load_keystore_directory()
    load_addresses()

    tests_results = []
    tests_results.append(test_mnemonics_from_sdk())

    for name in KEYS_ADDED:
        delete_from_keystore_by_name(name)

    if all(tests_results):
        print("\nPassed all tests!\n")
    else:
        print("\nFailed some tests, check logs.\n")
        sys.exit(-1)
