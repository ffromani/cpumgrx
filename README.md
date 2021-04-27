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
$ cpumgrx -M examples/machineinfo-v38.json -R 0,1 examples/trivial-pod.yaml 
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

example with non-guaranteed pods:
```bash
$ cpumgrx -v=5 -M examples/machineinfo-v38-lab.json -R "0,40,1,41" examples/nongu-pod.yaml 
I0427 12:45:53.411109   34156 cpu_manager.go:153] [cpumanager] detected CPU topology: &{104 52 2 map[0:{0 0 0} 1:{1 1 1} 2:{0 0 2} 3:{1 1 3} 4:{0 0 4} 5:{1 1 5} 6:{0 0 6} 7:{1 1 7} 8:{0 0 8} 9:{1 1 9} 10:{0 0 10} 11:{1 1 11} 12:{0 0 12} 13:{1 1 13} 14:{0 0 14} 15:{1 1 15} 16:{0 0 16} 17:{1 1 17} 18:{0 0 18} 19:{1 1 19} 20:{0 0 20} 21:{1 1 21} 22:{0 0 22} 23:{1 1 23} 24:{0 0 24} 25:{1 1 25} 26:{0 0 26} 27:{1 1 27} 28:{0 0 28} 29:{1 1 29} 30:{0 0 30} 31:{1 1 31} 32:{0 0 32} 33:{1 1 33} 34:{0 0 34} 35:{1 1 35} 36:{0 0 36} 37:{1 1 37} 38:{0 0 38} 39:{1 1 39} 40:{0 0 40} 41:{1 1 41} 42:{0 0 42} 43:{1 1 43} 44:{0 0 44} 45:{1 1 45} 46:{0 0 46} 47:{1 1 47} 48:{0 0 48} 49:{1 1 49} 50:{0 0 50} 51:{1 1 51} 52:{0 0 0} 53:{1 1 1} 54:{0 0 2} 55:{1 1 3} 56:{0 0 4} 57:{1 1 5} 58:{0 0 6} 59:{1 1 7} 60:{0 0 8} 61:{1 1 9} 62:{0 0 10} 63:{1 1 11} 64:{0 0 12} 65:{1 1 13} 66:{0 0 14} 67:{1 1 15} 68:{0 0 16} 69:{1 1 17} 70:{0 0 18} 71:{1 1 19} 72:{0 0 20} 73:{1 1 21} 74:{0 0 22} 75:{1 1 23} 76:{0 0 24} 77:{1 1 25} 78:{0 0 26} 79:{1 1 27} 80:{0 0 28} 81:{1 1 29} 82:{0 0 30} 83:{1 1 31} 84:{0 0 32} 85:{1 1 33} 86:{0 0 34} 87:{1 1 35} 88:{0 0 36} 89:{1 1 37} 90:{0 0 38} 91:{1 1 39} 92:{0 0 40} 93:{1 1 41} 94:{0 0 42} 95:{1 1 43} 96:{0 0 44} 97:{1 1 45} 98:{0 0 46} 99:{1 1 47} 100:{0 0 48} 101:{1 1 49} 102:{0 0 50} 103:{1 1 51}]}
I0427 12:45:53.411749   34156 policy_static.go:110] [cpumanager] reserved 4 CPUs ("0-1,40-41") not available for exclusive assignment
I0427 12:45:53.411795   34156 cpu_manager.go:193] [cpumanager] starting with static policy
I0427 12:45:53.411805   34156 cpu_manager.go:194] [cpumanager] reconciling every 10m0s
I0427 12:45:53.411833   34156 state_mem.go:36] [cpumanager] initializing new in-memory state store
I0427 12:45:53.417191   34156 state_mem.go:88] [cpumanager] updated default cpuset: "0-103"
I0427 12:45:53.421525   34156 main.go:140] handling pod: {"kind":"Pod","apiVersion":"v1","metadata":{"name":"qos-demo","namespace":"qos-example","creationTimestamp":null},"spec":{"containers":[{"name":"qos-demo-ctr","image":"nginx","resources":{"requests":{"cpu":"1100m","memory":"2Gi"}}}]},"status":{}}
qos-demo: 0-103 -> [ 45=[45,97] 14=[14,66] 12=[12,64] 25=[25,77] 24=[24,76] 47=[47,99] 18=[18,70] 50=[50,102] 31=[31,83] 11=[11,63] 32=[32,84] 46=[46,98] 48=[48,100] 30=[30,82] 36=[36,88] 43=[43,95] 38=[38,90] 37=[37,89] 44=[44,96] 8=[8,60] 40=[40,92] 33=[33,85] 6=[6,58] 15=[15,67] 42=[42,94] 7=[7,59] 9=[9,61] 3=[3,55] 23=[23,75] 1=[1,53] 34=[34,86] 19=[19,71] 28=[28,80] 21=[21,73] 4=[4,56] 49=[49,101] 35=[35,87] 10=[10,62] 16=[16,68] 5=[5,57] 41=[41,93] 2=[2,54] 17=[17,69] 29=[29,81] 0=[0,52] 27=[27,79] 22=[22,74] 13=[13,65] 51=[51,103] 26=[26,78] 39=[39,91] 20=[20,72] ]
00 -> [reserved qos-demo] <---
01 -> [reserved qos-demo] <---
02 -> [qos-demo]
03 -> [qos-demo]
04 -> [qos-demo]
05 -> [qos-demo]
06 -> [qos-demo]
07 -> [qos-demo]
08 -> [qos-demo]
09 -> [qos-demo]
10 -> [qos-demo]
11 -> [qos-demo]
12 -> [qos-demo]
13 -> [qos-demo]
14 -> [qos-demo]
15 -> [qos-demo]
16 -> [qos-demo]
17 -> [qos-demo]
18 -> [qos-demo]
19 -> [qos-demo]
20 -> [qos-demo]
21 -> [qos-demo]
22 -> [qos-demo]
23 -> [qos-demo]
24 -> [qos-demo]
25 -> [qos-demo]
26 -> [qos-demo]
27 -> [qos-demo]
28 -> [qos-demo]
29 -> [qos-demo]
30 -> [qos-demo]
31 -> [qos-demo]
32 -> [qos-demo]
33 -> [qos-demo]
34 -> [qos-demo]
35 -> [qos-demo]
36 -> [qos-demo]
37 -> [qos-demo]
38 -> [qos-demo]
39 -> [qos-demo]
40 -> [reserved qos-demo] <---
41 -> [reserved qos-demo] <---
42 -> [qos-demo]
43 -> [qos-demo]
44 -> [qos-demo]
45 -> [qos-demo]
46 -> [qos-demo]
47 -> [qos-demo]
48 -> [qos-demo]
49 -> [qos-demo]
50 -> [qos-demo]
51 -> [qos-demo]
I0427 12:45:53.424220   34156 main.go:122] removing cpu_manager_state file on "cpu_manager_state"
```

The format is `name=REQUEST/LIMIT`. You can use any valid format (see kubernetes' docs) for the resource spec.

## Obtaining machineinfos

1. [run cadvisor](https://github.com/google/cadvisor#quick-start-running-cadvisor-in-a-docker-container) on the box you want to collect the machineinfo for.
2. query the cadvisor API: `curl -L 127.0.0.1:8080/api/v1.3/machine > machineinfo.json`
3. feed the collected `machineinfo.json` into `cpumgrx`: `$ cpumgrx -M machineinfo.json -R 0,1 -T 'test1=1/1'`

Alternatively:
1. run `lsmachineinfo > machineinfo.json`

## license
(C) 2021 Red Hat Inc and licensed under the Apache License v2

## build
just run
```bash
make
```

## see also
similar tool to inspect topology manager behaviour: https://github.com/fromanirh/tmpolx
