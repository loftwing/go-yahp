:: move to vc dir
pushd "C:\Program Files (x86)\Microsoft Visual Studio\2017\Community\VC\Auxiliary\Build"

:: configure dev prompt
call .\vcvarsall.bat x64

:: back to project
popd

:: clean
del .\port.exe
:: compile: stack cookies, CF guard, sdl warnings/errors, only compile
cl /GS /guard:cf /sdl /c .\src\port.c

:: link: ASLR, CF guard
link /DYNAMICBASE /GUARD:cf /OUT:port.exe port.obj

:: cleanup
del port.obj

:: done
pause
exit
