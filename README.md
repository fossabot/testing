[![Build](https://github.com/tkrop/testing/actions/workflows/go.yaml/badge.svg)](https://github.com/tkrop/testing/actions/workflows/go.yaml)
[![Coverage](https://coveralls.io/repos/github/tkrop/testing/badge.svg?branch=main)](https://coveralls.io/github/tkrop/testing?branch=main)
[![Libraries](https://img.shields.io/hackage-deps/v/github.com/tkrop/testing)](https://img.shields.io/hackage-deps/v/github.com/tkrop/testing)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![FOSSA](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ftkrop%2Ftesting.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Ftkrop%2Ftesting?ref=badge_shield)
[![Report](https://goreportcard.com/badge/github.com/tkrop/testing)](https://goreportcard.com/badge/github.com/tkrop/testing)
[![Docs](https://pkg.go.dev/badge/github.com/tkrop/testing.svg)](https://pkg.go.dev/github.com/tkrop/testing)


# Testing

The testing framework contains a couple of small opinionated extensions for
[Golang][go], [Gomock][gomock], and [Gock][gock] to enable isolated parallel
parameterized unit and component tests using a common unified pattern to setup
mock requests and responses chains that work accross detached go routines.

```go
var testUnitCallParams = map[string]struct {
    mockSetup    mock.SetupFunc
    input*...    *model.*
    expect       test.Expect
    expect*...   *model.*
    expectError  error
}{
    "success" {
        mockSetup: mock.Chain(
            CallA(input..., output...), ...
        ),
        ...
        expect: test.ExpectSuccess
    }
}

func TestUnitCall(t *testing.T) {
    t.Parallel()

    for message, param := range testParams {
        // ensures copying parameters
        message, param := message, param
        t.Run(message, test.Run(param.expect, func(t test.Test) {
            t.Parallel()

            // Given
            mocks := mock.NewMock(t).Expect(param.mockSetup)

            unit := NewUnitService(
                mock.Get(mocks, NewServiceMock), ...
            )

            // When
            result, err := unit.call(param.input*...)

            mocks.Wait()

            // Then
            if param.expectError != nil {
                assert.Equal(t, param.expectError, err)
            } else {
                require.NoError(t, err)
            }
            assert.Equal(t, param.expect*, result)
        }))
    }
}
```


# Parallel test requirements

Running tests in parallel not only makes test faster, but also helps to detect
race conditions that else randomly appear in production  when running tests
with `go test -race`.

**Note:** there are some general requirements for running test in parallel:

1. Tests *must not modify* environment variables dynamically.
2. Tests *must not require* reserved service ports and open listeners.
3. Tests *must not share* resources, e.g. objects or database schemas, that
   are updated during execution of any parallel test.
4. Tests *must not use* [monkey patching][monkey] to modify commonly used
   functions, e.g. `time.Now()`, and finally
5. Tests *must not use* [Gock][gock] for mocking HTTP responses on transport
   level, instead use the [gock](gock)-controller provided by this framework.

If this conditions hold, the general pattern provided above can be used to
support parallel test execution.


# Project packages

The framework consists of the following sub-packages:

* [mock](mock) provides the means to setup a simple chain or a complex network
  of expected mock calls with minimal effort. This makes it easy to extend the
  usual narrow range of mocking to larger components using a unified pattern.

* [test](test) provides a small framework to simply isolate the test execution
  and safely check whether a test fails as expected. This is primarily very
  handy to validate a test framework as provided by the [mock](mock) package
  but may be handy in other cases too.

* [gock](gock) provides a drop-in extension for [Gock][gock] consisting of a
  controller and a mock storage that allows to run tests isolated. This allows
  to parallelize simple test and parameterized tests.

* [perm](perm) provides a small framework to simplify permutation tests, i.e.
  a consistent test set where conditions can be checked in all known orders
  with different outcome. This is very handy in combination with [test](test)
  to validated the [mock](mock) framework, but may be useful in other cases
  too.

Please see the documentation of the sub-packages for more details.


# Terms of Usage

This software is open source as is under the MIT license. If you start using
the software, please give it a star, so that I know to be more careful with
changes. If this project has more than 25 Stars, I will introduce semantic
versioning for changes.


# Contributing

If you like to contribute, please create an issue and/or pull request with a
proper description of your proposal or contribution. I will review it and
provide feedback on it.

[go]: https://go.dev/ "Golang"
[gomock]: https://github.com/golang/mock "GoMock"
[gock]: https://github.com/h2non/gock "Gock"
[monkey]: https://github.com/bouk/monkey "Monkey Patching"
