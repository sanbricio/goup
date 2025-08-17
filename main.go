package main

import (
	"fmt"
	"os"
	"os/exec"

	"golang.org/x/mod/modfile"
)

func main() {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error leyendo go.mod: %v\n", err)
		os.Exit(1)
	}
	f, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parseando go.mod: %v\n", err)
		os.Exit(1)
	}

	anyError := false

	for _, req := range f.Require {
		if !req.Indirect {
			fmt.Printf("Actualizando %s...\n", req.Mod.Path)
			cmd := exec.Command("go", "get", "-u", req.Mod.Path)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error actualizando %s: %v\n", req.Mod.Path, err)
				anyError = true
			}
		}
	}

	fmt.Println("Ejecutando go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	if anyError {
		fmt.Fprintln(os.Stderr, "¡Terminado con errores en alguna dependencia!")
		os.Exit(1)
	}
	fmt.Println("¡Actualización de dependencias directas terminada!")
}
