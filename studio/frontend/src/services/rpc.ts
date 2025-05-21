const self = globalThis as any;

type RPC = (...args: any[]) => any;

function isRPC(x: any): x is RPC {
  return typeof x === "function";
}

export async function call(name: string, ...args: any[]) {
  const caller = self["rpc_" + name];
  if (isRPC(caller)) {
    console.info("calling", name, args);
    const output = await caller(...args);
    console.info("output", output);
    return output;
  }
  throw new Error(`RPC ${name} not found`);
}
