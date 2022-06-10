import os
import shutil
import time
import sys

BIN_DIR = os.path.join(os.path.expanduser("~"), "bin")

bashrc = os.path.join(os.path.expanduser("~"), ".bashrc")
bashProfile = os.path.join(os.path.expanduser("~"), ".bash_profile")
zshRc = os.path.join(os.path.expanduser("~"), ".zshrc")

if len(sys.argv) >= 2:
    if sys.argv[1] == "--all":
        print("Installing to all RCs")
        for rc in [bashrc, bashProfile, zshRc]:
            if not os.path.exists(rc):
                print("Creating file: " + rc)
                with open(rc, "w") as f:
                    f.write("")

print("Detecting .rc files...")

foundFiles = []

if os.path.isfile(bashrc):
    print("Found .bashrc")
    foundFiles.append(bashrc)

if os.path.isfile(bashProfile):
    print("Found .bash_profile")
    foundFiles.append(bashProfile)

if os.path.isfile(zshRc):
    print("Found .zshrc")
    foundFiles.append(zshRc)

if len(foundFiles) == 0:
    print("No .rc files found, creating bashrc, bash_profile and zshrc")
    with open(bashrc, "w") as f:
        f.write("\n")
    with open(bashProfile, "w") as f:
        f.write("\n")
    with open(zshRc, "w") as f:
        f.write("\n")
    foundFiles.extend([bashrc, bashProfile, zshRc])

print("Adding the bin directory to the rc files...")
for filePath in foundFiles:
    print("Backing up " + filePath)
    shutil.copy(filePath, f"{filePath}.{str(int(time.time()))}.yyjlincoln.bak")
    with open(filePath, "r") as file:
        if file.read().find("export PATH=\"" + BIN_DIR + ":$PATH\"\n") != -1:
            print("Already added to " + filePath)
            continue
    with open(filePath, "a") as f:
        f.write("\n")
        f.write("export PATH=\"" + BIN_DIR + ":$PATH\"\n")
    print("Added to " + filePath)

print("Installation complete.")
