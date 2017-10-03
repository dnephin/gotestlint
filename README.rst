gotestlint
==========

``gotestlint`` checks that all tests conform to the following convention:

* the testcase is named in the form ``Test<Function>[<Condition>]``
* the testcase tests a function in the same pacakge
* the testcase is in a file name that mirrors the filename of the function
  under test, with a ``_test`` suffix.
