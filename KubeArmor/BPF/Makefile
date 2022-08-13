RED=\033[0;31m
GREEN=\033[0;32m
CYAN=\033[0;36m
NC=\033[0m

KRBTF = /sys/kernel/btf/vmlinux

KRNVER = $(shell uname -r)
# split kernel version string
KRNV_S = $(subst ., ,$(KRNVER)) # x.y.z-rc-arch --> x y z-rc-arch
KRNV_X = $(word 1,$(KRNV_S)) # major version
KRNV_Y = $(word 2,$(KRNV_S)) # patch level version 
KRNV_Z = $(firstword $(subst -, ,$(word 3,$(KRNV_S)))) # sublevel version

ifeq (,$(KRNDIR))
ifneq (,$(wildcard /lib/modules/$(KRNVER)/build/Kconfig))
	KRNDIR = /lib/modules/$(KRNVER)/build
else ifneq (,$(wildcard /lib/modules/$(KRNVER)/source/Kconfig))
	KRNDIR = /lib/modules/$(KRNVER)/source
else ifneq (,$(wildcard /etc/redhat-release))
	KRNDIR = /usr/src/$(KRNVER)
else
	KRNDIR = /usr/src/linux-headers-$(KRNVER)
endif
endif

LIBBPF = $(CURDIR)/libbpf

CL  = clang
CC  = gcc
Q   = @

ifeq ($(V),1)
  Q =
endif

BTF_SUPPORTED = 0
ifneq (,$(wildcard $(KRBTF)))
	BTF_SUPPORTED = 1
endif

ifeq ($(BTF_SUPPORTED),1)
	KF = -DBTF_SUPPORTED
else
	# copied from kernel's samples/bpf/Makefile
	KF = -I$(KRNDIR)/arch/x86/include -I$(KRNDIR)/arch/x86/include/generated  \
	 	 -I$(KRNDIR)/include -I$(KRNDIR)/arch/x86/include/uapi \
	 	 -I$(KRNDIR)/arch/x86/include/generated/uapi -I$(KRNDIR)/include/uapi \
	 	 -I$(KRNDIR)/include/generated/uapi \
	 	 -include $(KRNDIR)/include/linux/kconfig.h
endif

KF := $(KF) -I$(LIBBPF)/src \
	 -DLINUX_VERSION_MAJOR=$(KRNV_X) \
	 -DLINUX_VERSION_PATCHLEVEL=$(KRNV_Y) \
	 -DLINUX_VERSION_SUBLEVEL=$(KRNV_Z) \
	 -D__KERNEL__ \
	 -D__BPF_TRACING__ \
	 -D__TARGET_ARCH_x86 \
	 -Wunused \
	 -Wno-frame-address \
	 -Wno-unused-value \
	 -Wno-unused-function \
	 -Wno-unknown-warning-option \
	 -Wno-pragma-once-outside-header \
	 -Wno-pointer-sign \
	 -Wno-gnu-variable-sized-type-not-at-end \
	 -Wno-deprecated-declarations \
	 -Wno-compare-distinct-pointer-types \
	 -Wno-address-of-packed-member \
	 -fno-stack-protector \
	 -fno-jump-tables \
	 -fno-unwind-tables \
	 -fno-asynchronous-unwind-tables \
	 -xc -O2 -g -target bpf

SYSMONITOR = $(CURDIR)/system_monitor.c
 
.PHONY: all
all: kernel_headers libbpf system_monitor.bpf.o system_monitor.container.bpf.o system_monitor.host.bpf.o

.PHONY: kernel_headers
kernel_headers:
ifeq ($(BTF_SUPPORTED),1)
	@echo "$(GREEN)Kernel BTF information found$(NC)"
	@echo "Generating vmlinux.h.."
	$(Q)bpftool btf dump file /sys/kernel/btf/vmlinux format c > $(CURDIR)/vmlinux.h
else
	@echo "$(RED)Kernel BTF information not found$(NC)"
	@echo "Trying to use kernel headers present in the host"
	ifeq (,$(wildcard $(KRNDIR)/Kconfig))
		@echo "$(RED)Auto-detected kernel headers path [$(KRNDIR)] is incorrect.$(NC) Use 'make KRNDIR=[KERNEL-SRC-PATH]'."
		Quitting 
	else
		@echo "$(GREEN)Using kernel headers at $(KRNDIR)$(NC)"
	endif
endif

.PHONY: libbpf
libbpf:
ifeq (,$(wildcard $(LIBBPF)/src/libbpf.c))
	$(Q)git submodule update --init --recursive
endif
ifeq (,$(wildcard $(LIBBPF)/src/libbpf.a))
	$(Q)make -C $(LIBBPF)/src
endif

system_monitor.bpf.o: $(SYSMONITOR)
	@echo "Compiling eBPF bytecode: $(GREEN)$@$(NC) ..."
	$(Q)$(CL) $(KF) -DMONITOR_HOST_AND_CONTAINER -c $< -o $@

system_monitor.container.bpf.o: $(SYSMONITOR)
	@echo "Compiling eBPF bytecode: $(GREEN)$@$(NC) ..."
	$(Q)$(CL) $(KF) -DMONITOR_CONTAINER -c $< -o $@

system_monitor.host.bpf.o: $(SYSMONITOR)
	@echo "Compiling eBPF bytecode: $(GREEN)$@$(NC) ..."
	$(Q)$(CL) $(KF) -DMONITOR_HOST -c $< -o $@

# clean up
clean:
	$(Q)rm -rf *.o vmlinux.h
	$(Q)make clean -C $(LIBBPF)/src
