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

cpumanager cares only about resource requirements. So if you don't want to supply a full pod spec every time, you can make `cpumgrx` autogenerate
a minimal fake pod spec, using the template mode:
```bash
$ cpumgrx -M examples/machineinfo-v38.json -R 0,1 -T 'test1=2/2' 'test2=5/5' 'test3=1/1' 'test4=8/8'
I0212 12:31:33.824419   43643 cpu_manager.go:153] [cpumanager] detected CPU topology: &{24 12 2 map[0:{0 1 0} 1:{1 0 1} 2:{0 1 2} 3:{1 0 3} 4:{0 1 4} 5:{1 0 5} 6:{0 1 6} 7:{1 0 7} 8:{0 1 8} 9:{1 0 9} 10:{0 1 10} 11:{1 0 11} 12:{0 1 0} 13:{1 0 1} 14:{0 1 2} 15:{1 0 3} 16:{0 1 4} 17:{1 0 5} 18:{0 1 6} 19:{1 0 7} 20:{0 1 8} 21:{1 0 9} 22:{0 1 10} 23:{1 0 11}]}
I0212 12:31:33.824808   43643 policy_static.go:110] [cpumanager] reserved 2 CPUs ("0-1") not available for exclusive assignment
I0212 12:31:33.824838   43643 cpu_manager.go:193] [cpumanager] starting with static policy
I0212 12:31:33.824846   43643 cpu_manager.go:194] [cpumanager] reconciling every 10m0s
I0212 12:31:33.824869   43643 state_mem.go:36] [cpumanager] initializing new in-memory state store
I0212 12:31:33.834711   43643 state_mem.go:88] [cpumanager] updated default cpuset: "0-23"
I0212 12:31:33.842316   43643 policy_static.go:221] [cpumanager] static policy: Allocate (pod: test1-pod_(), container: test1-cnt)
I0212 12:31:33.842353   43643 policy_static.go:232] [cpumanager] Pod test1-pod_(), Container test1-cnt Topology Affinity is: {<nil> false}
I0212 12:31:33.842377   43643 policy_static.go:259] [cpumanager] allocateCpus: (numCPUs: 2, socket: <nil>)
I0212 12:31:33.842623   43643 state_mem.go:88] [cpumanager] updated default cpuset: "0-2,4-14,16-23"
I0212 12:31:33.852995   43643 policy_static.go:294] [cpumanager] allocateCPUs: returning "3,15"
I0212 12:31:33.853068   43643 state_mem.go:80] [cpumanager] updated desired cpuset (pod: , container: test1-cnt, cpuset: "3,15")
3,15
I0212 12:31:33.860657   43643 policy_static.go:221] [cpumanager] static policy: Allocate (pod: test2-pod_(), container: test2-cnt)
I0212 12:31:33.860714   43643 policy_static.go:232] [cpumanager] Pod test2-pod_(), Container test2-cnt Topology Affinity is: {<nil> false}
I0212 12:31:33.860747   43643 policy_static.go:259] [cpumanager] allocateCpus: (numCPUs: 5, socket: <nil>)
I0212 12:31:33.862237   43643 state_mem.go:88] [cpumanager] updated default cpuset: "0-2,4,6,8-11,13-14,16,18,20-23"
I0212 12:31:33.869182   43643 policy_static.go:294] [cpumanager] allocateCPUs: returning "5,7,12,17,19"
I0212 12:31:33.869260   43643 state_mem.go:80] [cpumanager] updated desired cpuset (pod: , container: test2-cnt, cpuset: "5,7,12,17,19")
5,7,12,17,19
I0212 12:31:33.876346   43643 policy_static.go:221] [cpumanager] static policy: Allocate (pod: test3-pod_(), container: test3-cnt)
I0212 12:31:33.876379   43643 policy_static.go:232] [cpumanager] Pod test3-pod_(), Container test3-cnt Topology Affinity is: {<nil> false}
I0212 12:31:33.876412   43643 policy_static.go:259] [cpumanager] allocateCpus: (numCPUs: 1, socket: <nil>)
I0212 12:31:33.878175   43643 state_mem.go:88] [cpumanager] updated default cpuset: "0-2,4,6,8-11,14,16,18,20-23"
I0212 12:31:33.882089   43643 policy_static.go:294] [cpumanager] allocateCPUs: returning "13"
I0212 12:31:33.882149   43643 state_mem.go:80] [cpumanager] updated desired cpuset (pod: , container: test3-cnt, cpuset: "13")
13
I0212 12:31:33.886251   43643 policy_static.go:221] [cpumanager] static policy: Allocate (pod: test4-pod_(), container: test4-cnt)
I0212 12:31:33.886318   43643 policy_static.go:232] [cpumanager] Pod test4-pod_(), Container test4-cnt Topology Affinity is: {<nil> false}
I0212 12:31:33.886366   43643 policy_static.go:259] [cpumanager] allocateCpus: (numCPUs: 8, socket: <nil>)
I0212 12:31:33.886655   43643 state_mem.go:88] [cpumanager] updated default cpuset: "0-1,6,8,10,18,20,22"
I0212 12:31:33.890802   43643 policy_static.go:294] [cpumanager] allocateCPUs: returning "2,4,9,11,14,16,21,23"
I0212 12:31:33.890877   43643 state_mem.go:80] [cpumanager] updated desired cpuset (pod: , container: test4-cnt, cpuset: "2,4,9,11,14,16,21,23")
2,4,9,11,14,16,21,23
```

The format is `name=REQUEST/LIMIT`. You can use any valid format (see kubernetes' docs) for the resource spec.

## Obtaining machineinfos

1. [run cadvisor](https://github.com/google/cadvisor#quick-start-running-cadvisor-in-a-docker-container) on the box you want to collect the machineinfo for.
2. query the cadvisor API: `curl -L 127.0.0.1:8080/api/v1.3/machine > machineinfo.json`
3. feed the collected `machineinfo.json` into `cpumgrx`: `$ cpumgrx -M machineinfo.json -R 0,1 -T 'test1=1/1'`

## license
(C) 2021 Red Hat Inc and licensed under the Apache License v2

## build
just run
```bash
make
```

## see also
similar tool to inspect topology manager behaviour: https://github.com/fromanirh/tmpolx
