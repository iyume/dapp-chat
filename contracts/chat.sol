// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Chat {
    uint value;

    function increase(uint a) public {
        value += a;
    }

    function getValue() public view returns (uint) {
        return value;
    }
}
