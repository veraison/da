SRCS := $(wildcard *.diag)
CBOR := $(SRCS:.diag=.cbor)
HEX := $(SRCS:.diag=.hex)

%.cbor: %.diag ; diag2cbor.rb < $< > $@

CLEANFILES += $(CBOR)

%.hex: %.cbor ; xxd -p -c0 $< > $@

CLEANFILES += $(HEX)

all: $(CBOR) $(B64U) $(HEX)

clean: ; $(RM) -f $(CLEANFILES)
