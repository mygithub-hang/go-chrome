package go_chrome

import (
	"github.com/akavel/rsrc/rsrc"
)

var manifestContent string = `
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
<assemblyIdentity
    version="1.0.0.0"
    processorArchitecture="x86"
    name="controls"
    type="win32"
></assemblyIdentity>
<dependency>
    <dependentAssembly>
        <assemblyIdentity
            type="win32"
            name="Microsoft.Windows.Common-Controls"
            version="6.0.0.0"
            processorArchitecture="*"
            publicKeyToken="6595b64144ccf1df"
            language="*"
        ></assemblyIdentity>
    </dependentAssembly>
</dependency>
</assembly>`

func BuildSysO() error {
	var fnamein, fnameico, fnameout, arch string
	// rsrc -manifest main.manifest -ico resources/icons/icon.ico -o main.syso
	fnamein = "main.manifest"
	fnameico = "resources/icons/icon.ico"
	fnameout = "main.syso"
	arch = "amd64"
	if fnameout == "" {
		fnameout = "rsrc_windows_" + arch + ".syso"
	}
	err := rsrc.Embed(fnameout, arch, fnamein, fnameico)
	return err
}
