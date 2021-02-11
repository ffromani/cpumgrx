# cpumgrx: CPU ManaGeR policy eXploration tool

A simple tool to test how kubernetes' CPU manager behave, without need to run the real workload.
This tool uses the very same packages from upstream kubernetes, to give the closest representation as
possible as the real thing.

Usage example:
```bash
$ cat examples/trivial-pod.yaml 
apiVersion: v1
kind: Pod
metadata:
  name: qos-demo
  namespace: qos-example
spec:
  containers:
  - name: qos-demo-ctr
    image: nginx
    resources:
      limits:
        memory: "2048Mi"
        cpu: "5"
      requests:
        memory: "2048Mi"
        cpu: "5"
$ cpumgrx -M examples/machineinfo-v38.json -R 0,1 -p examples/trivial-pod.yaml 
I0211 15:26:38.403129   66576 cpu_manager.go:153] [cpumanager] detected CPU topology: &{24 12 2 map[0:{0 1 0} 1:{1 0 1} 2:{0 1 2} 3:{1 0 3} 4:{0 1 4} 5:{1 0 5} 6:{0 1 6} 7:{1 0 7} 8:{0 1 8} 9:{1 0 9} 10:{0 1 10} 11:{1 0 11} 12:{0 1 0} 13:{1 0 1} 14:{0 1 2} 15:{1 0 3} 16:{0 1 4} 17:{1 0 5} 18:{0 1 6} 19:{1 0 7} 20:{0 1 8} 21:{1 0 9} 22:{0 1 10} 23:{1 0 11}]}
I0211 15:26:38.403559   66576 policy_static.go:110] [cpumanager] reserved 2 CPUs ("0-1") not available for exclusive assignment
I0211 15:26:38.403591   66576 cpu_manager.go:193] [cpumanager] starting with static policy
I0211 15:26:38.403604   66576 cpu_manager.go:194] [cpumanager] reconciling every 10m0s
I0211 15:26:38.403631   66576 state_mem.go:36] [cpumanager] initializing new in-memory state store
I0211 15:26:38.420395   66576 state_mem.go:88] [cpumanager] updated default cpuset: "0-23"
I0211 15:26:38.432347   66576 policy_static.go:221] [cpumanager] static policy: Allocate (pod: qos-demo_qos-example(), container: qos-demo-ctr)
I0211 15:26:38.432932   66576 policy_static.go:232] [cpumanager] Pod qos-demo_qos-example(), Container qos-demo-ctr Topology Affinity is: {<nil> false}
I0211 15:26:38.432957   66576 policy_static.go:259] [cpumanager] allocateCpus: (numCPUs: 5, socket: <nil>)
I0211 15:26:38.433519   66576 state_mem.go:88] [cpumanager] updated default cpuset: "0-2,4,6-12,14,16,18-23"
I0211 15:26:38.437875   66576 policy_static.go:294] [cpumanager] allocateCPUs: returning "3,5,13,15,17"
I0211 15:26:38.437930   66576 state_mem.go:80] [cpumanager] updated desired cpuset (pod: , container: qos-demo-ctr, cpuset: "3,5,13,15,17")
3,5,13,15,17
```

## license
(C) 2021 Red Hat Inc and licensed under the Apache License v2

## build
just run
```bash
make
```
