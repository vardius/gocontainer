// Results behind the benchmark charts and tables, updated together. The suite reruns with the commands on the
// Benchmarks page; the comparison comes from a harness that runs the old and new versions side by side.

export const environment = 'Apple M4 Max, Go 1.27.0, median of 5 runs';

/** The two versions the charts compare. */
export const versions = { before: 'v1', after: 'v2' } as const;

export type Group = 'get' | 'invoke' | 'parallel' | 'register';

export interface Result {
	ns: number;
	allocs: number;
}

export interface Comparison {
	scenario: string;
	before: Result;
	after: Result;
}

/**
 * v1 against v2 on the same harness, one call per op, in a container of 32 services. v1's Invoke calls stand
 * against the v2 calls that replace them.
 */
export const comparison: Record<Group, Comparison[]> = {
	get: [
		{ scenario: 'Get', before: { ns: 5.8, allocs: 0 }, after: { ns: 3.6, allocs: 0 } },
		{ scenario: 'Get, not registered', before: { ns: 4.8, allocs: 0 }, after: { ns: 3.3, allocs: 0 } },
		{ scenario: 'MustGet', before: { ns: 6.4, allocs: 0 }, after: { ns: 3.5, allocs: 0 } },
		{ scenario: 'Has', before: { ns: 6.2, allocs: 0 }, after: { ns: 3.4, allocs: 0 } },
	],
	invoke: [
		{ scenario: 'Invoke in v1, Get in v2', before: { ns: 80.9, allocs: 0 }, after: { ns: 3.6, allocs: 0 } },
		{ scenario: 'MustInvoke in v1, MustGet in v2', before: { ns: 71, allocs: 0 }, after: { ns: 3.5, allocs: 0 } },
		{
			scenario: 'Three services: MustInvokeMany in v1, three MustGets in v2',
			before: { ns: 162.7, allocs: 5 },
			after: { ns: 10.2, allocs: 0 },
		},
	],
	parallel: [
		{ scenario: 'Get from 16 goroutines at once', before: { ns: 128.1, allocs: 0 }, after: { ns: 0.27, allocs: 0 } },
	],
	register: [
		{ scenario: 'Register, replacing a service', before: { ns: 9.5, allocs: 0 }, after: { ns: 224.6, allocs: 5 } },
	],
};

export interface SuiteResult {
	benchmark: string;
	what: string;
	ns: number;
	bytes: number;
	allocs: number;
}

/** v2's own benchmark suite, `go test -bench .`, in a container of 32 services unless the name says otherwise. */
export const suite: SuiteResult[] = [
	{ benchmark: 'Get', what: 'Get a registered service', ns: 3.5, bytes: 0, allocs: 0 },
	{ benchmark: 'GetNotRegistered', what: 'Get under a key with nothing registered', ns: 3.4, bytes: 0, allocs: 0 },
	{ benchmark: 'MustGet', what: 'MustGet a registered service', ns: 3.4, bytes: 0, allocs: 0 },
	{ benchmark: 'Has', what: 'Check for a registered service', ns: 3.3, bytes: 0, allocs: 0 },
	{ benchmark: 'GetParallel', what: 'Get from 16 goroutines at once', ns: 0.26, bytes: 0, allocs: 0 },
	{ benchmark: 'Register/1', what: 'Replace the only service', ns: 52.7, bytes: 264, allocs: 3 },
	{ benchmark: 'Register/32', what: 'Replace one of 32 services', ns: 216.8, bytes: 1888, allocs: 5 },
	{ benchmark: 'Register/1024', what: 'Replace one of 1024 services', ns: 4710, bytes: 54664, allocs: 7 },
];
