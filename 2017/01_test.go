package main

import (
    "testing"
)

func TestSum(t *testing.T){
    testCases := []struct{
        have string
        want int
    }{
        { "1122", 3 },
        { "1111", 4 },
        { "1234", 0 },
        { "91212129", 9 },
    }

     for _, tc := range testCases {
        got := getSum(tc.have)
        if got != tc.want {
            t.Fatalf("Have %v, want %v, got %v", tc.have, tc.want, got)
        }
    }
}