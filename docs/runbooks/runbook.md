---
title: SRE Runbook
url: https://github.com/matthisholleville/ava/blob/main/docs/runbooks/runbook.md
---


# SRE Runbook


## Pod Crashlooping


In case of a crashloop of one or more pods, this may impact the application's availability. Please follow the procedure:

1. check if http://198.19.249.2/ is available
2. Retrieve the details of the pod that is in CrashLoopBackOff state.
3. Collect the logs of the pod.
4. Delete the affected pod.
5. Wait 10 seconds
5. Check if http://198.19.249.2/ is available