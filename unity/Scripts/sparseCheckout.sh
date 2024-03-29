#!/bin/sh
echo "this file can be put into the /Assets/unity/Scripts/sparseCheckout.sh"
cd $(dirname $0)
pwd
git init
git remote add origin git@github.com-pivot:GitVerse-Web3/PIVOT.git
git config core.sparseCheckout true
echo "unity/Scripts/" >> .git/info/sparse-checkout
git pull origin master
