#!/usr/bin/env python3

import os
import json
from subprocess import run, STDOUT, PIPE, DEVNULL

import gi

gi.require_version("Gtk", "3.0")

from gi.repository import Gio, Gtk

tmp_file = "/tmp/launcher.json"


def icon_imgcat(icon_path):
    p = run(["wezterm", "imgcat", "--height", "1", icon_path], capture_output=True)

    if p.returncode != 0:
        return ""
    else:
        s = p.stdout.decode()
        if s.endswith("\n"):
            s = s[:-1]
        return s


def load_applications():
    with open(tmp_file, "r") as f:
        applications = json.load(f)

    return applications


def on_open(*_, **__):
    gio_apps = Gio.AppInfo.get_all()
    icon_theme = Gtk.IconTheme.get_default()

    applications = []

    # print("Loaded applications:")
    for gio_app in gio_apps:
        # print(gio_app.get_name())
        # print(gio_app.get_filename())
        # print(gio_app.get_description())

        app = {
            "id": gio_app.get_id(),
            "name": gio_app.get_name(),
            "filename": gio_app.get_filename(),
            "description": gio_app.get_description(),
        }

        icon = gio_app.get_icon()
        if icon:
            # print(gio_app.get_icon().to_string())

            icon_info = icon_theme.lookup_icon(icon.to_string(), 48, 0)

            if icon_info:
                # print(icon_info.get_filename())
                app["icon"] = icon_info.get_filename()
            # else:
            #     print("No icon found")

        # print()

        applications.append(app)

    with open(tmp_file, "w") as f:
        json.dump(applications, f)

    on_draw(applications, prompt="")


def on_draw(apps, prompt, *_, **__):
    lines = []

    for app in apps:
        # Icons do not render correctly in the menu and are very slow

        # icon = icon_imgcat(app["icon"]) if "icon" in app else "   "
        # lines.append(f"{icon} {app["name"]} | {app["description"]} | {app['id']}")

        lines.append(f"{app["name"]} | {app["description"]} | {app['id']}")

    p = run(["menu", "search", prompt], input="\n".join(lines), text=True, stdout=PIPE, stderr=DEVNULL)
    if p.returncode != 0:
        print("No results")
    else:
        print(p.stdout)


def on_close(apps, sel_line, *_, **__):
    sel_id = sel_line.split(" | ")[-1]
    apps_by_id = {app["id"]: app for app in apps}
    sel_app = apps_by_id.get(sel_id)

    run(["gio", "launch", sel_app["filename"]])
    exit()


handlers = {"open": on_open, "close": on_close}


def main():
    prompt = os.getenv("prompt")
    event = os.getenv("event")
    sel_line = os.getenv("sel_line")

    apps = []
    if event != "open":
        apps = load_applications()

    handlers.get(event, on_draw)(apps, prompt=prompt, sel_line=sel_line)


if __name__ == "__main__":
    main()
