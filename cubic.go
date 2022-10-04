// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
    "fmt"
    "os"
    "github.com/consensys/gnark/frontend"
    "github.com/consensys/gnark-crypto/ecc"
    "github.com/consensys/gnark/frontend/cs/r1cs"
    "github.com/consensys/gnark/backend/groth16"
    _ "gnark-ed25519/edwards_curve"
    _ "gnark-ed25519/sha512"
)

// Circuit defines a simple circuit
// x**3 + x + 5 == y
type Circuit struct {
    // struct tags on a variable is optional
    // default uses variable name and secret visibility.
    X frontend.Variable `gnark:"x"`
    Y frontend.Variable `gnark:",public"`
}

// Define declares the circuit constraints
// x**3 + x + 5 == y
func (circuit *Circuit) Define(api frontend.API) error {
    x3 := api.Mul(circuit.X, circuit.X, circuit.X)
    api.AssertIsEqual(circuit.Y, api.Add(x3, circuit.X, 5))
    return nil
}

func main() {
    err := mainImpl()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func mainImpl() error {
    var myCircuit Circuit
    r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &myCircuit)
    if err != nil {
        return err
    }

    assignment := &Circuit{
        X: "2",
        Y: "15",
    }
    witness, _ := frontend.NewWitness(assignment, ecc.BN254.ScalarField())
    publicWitness, _ := witness.Public()
    pk, vk, err := groth16.Setup(r1cs)
    proof, err := groth16.Prove(r1cs, pk, witness)
    err = groth16.Verify(proof, vk, publicWitness)
    if err != nil {
        return err
    }
    fmt.Println(proof)
    return nil
}