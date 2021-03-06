# TINGen
TINGen creates German TINs (in German: "Steueridentifikationsnummer") like desribed [here](https://download.elster.de/download/schnittstellen/Pruefung_der_Steuer_und_Steueridentifikatsnummer.pdf) (German). It is also able to create TINs for testing purpose which could, in contrast to non-testing TINs, start with zero. See the [Wikipedia article](https://de.wikipedia.org/wiki/Steuerliche_Identifikationsnummer) about TIN (German).
## Screenshot
![TINGen screenshot](https://raw.githubusercontent.com/dartmann/TINGen/master/img/TINGen.PNG "Screenshot of TINGen UI")
## Building
1. Follow the steps described [here](https://github.com/fyne-io/fyne/wiki/Compiling) for being able to build TINGen (this is needed for the GUI).
2. Run `Go build` in the folder of TINGen.go file.
3. An TINGen.exe file should be compiled.
4. Profit.
## GUI
The graphical user interface is built with [Fyne](https://github.com/fyne-io/fyne).