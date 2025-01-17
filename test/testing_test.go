package test

import (
	"testing"

	"github.com/tkrop/testing/mock"
)

//go:generate mockgen -package=test -destination=mock_iface_test.go -source=testing_test.go  IFace

type IFace interface {
	CallA(string)
	CallB(string) string
}

func CallA(input string) mock.SetupFunc {
	return func(mocks *mock.Mocks) any {
		return mock.Get(mocks, NewMockIFace).EXPECT().
			CallA(input).Times(mocks.Times(1)).
			Do(mocks.GetDone(1))
	}
}

func CallB(input string, output string) mock.SetupFunc {
	return func(mocks *mock.Mocks) any {
		return mock.Get(mocks, NewMockIFace).EXPECT().
			CallB(input).Return(output).
			Times(mocks.Times(1)).Do(mocks.GetDone(1))
	}
}

var testRunParams = map[string]struct {
	expect Expect
	test   func(Test)
}{
	"run success": {
		test: InRun(ExpectSuccess,
			func(Test) {}),
		expect: ExpectSuccess,
	},

	"run success with errorf": {
		test: InRun(ExpectSuccess,
			func(t Test) { t.Errorf("fail") }),
		expect: ExpectFailure,
	},

	"run success with fatalf": {
		test: InRun(ExpectSuccess,
			func(t Test) { t.Fatalf("fail") }),
		expect: ExpectFailure,
	},

	"run success with failnow": {
		test: InRun(ExpectSuccess,
			func(t Test) { t.FailNow() }),
		expect: ExpectFailure,
	},

	"run failure": {
		test: InRun(ExpectFailure,
			func(t Test) {}),
		expect: ExpectFailure,
	},

	"run failure with errorf": {
		test: InRun(ExpectFailure,
			func(t Test) { t.Errorf("fail") }),
		expect: ExpectSuccess,
	},

	"run failure with fatalf": {
		test: InRun(ExpectFailure,
			func(t Test) { t.Fatalf("fail") }),
		expect: ExpectFailure,
	},

	"run failure with failnow": {
		test: InRun(ExpectFailure,
			func(t Test) { t.FailNow() }),
		expect: ExpectFailure,
	},

	"in success": {
		test:   InSuccess(func(t Test) {}),
		expect: ExpectSuccess,
	},

	"in success with errorf": {
		test: InSuccess(
			func(t Test) { t.Errorf("fail") }),
		expect: ExpectFailure,
	},

	"in success with fatalf": {
		test: InSuccess(
			func(t Test) { t.Fatalf("fail") }),
		expect: ExpectFailure,
	},

	"in success with failnow": {
		test: InSuccess(
			func(t Test) { t.FailNow() }),
		expect: ExpectFailure,
	},

	"in failure": {
		test:   InFailure(func(t Test) {}),
		expect: ExpectFailure,
	},

	"in failure with errorf": {
		test: InFailure(
			func(t Test) { t.Errorf("fail") }),
		expect: ExpectSuccess,
	},

	"in failure with fatalf": {
		test: InFailure(
			func(t Test) { t.Fatalf("fail") }),
		expect: ExpectFailure,
	},

	"in failure with failnow": {
		test: InFailure(
			func(t Test) { t.FailNow() }),
		expect: ExpectFailure,
	},
}

func Call(t Test, mocks *mock.Mocks, expect Expect, test func(Test)) {
	test(t)
	if expect == ExpectSuccess {
		mock.Get(mocks, NewMockIFace).CallA("a")
	} else {
		// simulate mock call since consumption of
		// mock call will creat a random result after
		// unlocking test thread.
		mocks.Times(-1)
	}
}

func TestRun(t *testing.T) {
	t.Parallel()

	for message, param := range testRunParams {
		message, param := message, param
		t.Run(message, Run(param.expect, func(t Test) {
			t.Parallel()

			// Given
			mocks := mock.NewMock(t).Expect(CallA("a"))

			// When
			go Call(t, mocks, param.expect, param.test)

			// Then
			mocks.Wait()
		}))
	}
}

func TestOther(t *testing.T) {
	t.Parallel()

	for message, param := range testRunParams {
		message, param := message, param
		switch param.expect {
		case ExpectSuccess:
			t.Run(message, Success(func(t Test) {
				t.Parallel()

				// Given
				mocks := mock.NewMock(t).Expect(CallA("a"))

				// When
				go Call(t, mocks, param.expect, param.test)

				// Then
				mocks.Wait()
			}))

		case ExpectFailure:
			t.Run(message, Failure(func(t Test) {
				t.Parallel()

				// Given
				mocks := mock.NewMock(t).Expect(CallA("a"))

				// When
				go Call(t, mocks, param.expect, param.test)

				// Then
				mocks.Wait()
			}))
		}
	}
}
