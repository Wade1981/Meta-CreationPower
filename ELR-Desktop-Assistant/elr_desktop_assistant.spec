# -*- mode: python ; coding: utf-8 -*-


a = Analysis(
    ['elr_desktop_assistant.py'],
    pathex=[],
    binaries=[],
    datas=[('icons/elr_icon.png', 'icons'), ('icons/elr_icon.ico', 'icons')],
    hiddenimports=[],
    hookspath=[],
    hooksconfig={},
    runtime_hooks=[],
    excludes=[],
    noarchive=False,
    optimize=0,
)
pyz = PYZ(a.pure)

exe = EXE(
    pyz,
    a.scripts,
    a.binaries,
    a.datas,
    [],
    name='elr_desktop_assistant',
    debug=False,
    bootloader_ignore_signals=False,
    strip=False,
    upx=True,
    upx_exclude=[],
    runtime_tmpdir=None,
    console=False,
    disable_windowed_traceback=False,
    argv_emulation=False,
    target_arch=None,
    codesign_identity=None,
    entitlements_file=None,
    icon=['E:\\X54\\github\\Meta-CreationPower\\ELR-Desktop-Assistant\\icons\\elr_icon.ico'],
)
