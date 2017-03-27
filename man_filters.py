#!/usr/bin/env python3

panflute_import = (
    "\n\n"
    "Cannot import `panflute` library, verify it's installed correctly\n"
    "For virtualenv:\n"
    "\tpip install panflute\n"
    "For system:\n"
    "\tpip install --user panflute\n"
)

try:
    from panflute import Code, Str, Strong, Image, run_filters
    # for debug purpose it's possible to use `debug`
    # function from tools module
    #
    # >>> from panflute.tools import debug

except ImportError:
    raise Exception(panflute_import)


def unwrap_code(elem, doc):
    if type(elem) == Code:
        text = elem.text.replace("bash", "")
        return [Str("\n"), Strong(Str(text))]

def remove_badge(elem, doc):
    if type(elem) == Image:
        return []

def main(doc=None):
    return run_filters([remove_badge, unwrap_code], doc=doc)

if __name__ == "__main__":
    main()
