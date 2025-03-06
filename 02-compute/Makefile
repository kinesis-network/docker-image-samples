OUTDIR=bin
OBJDIR=obj
SRCDIR=src

CC=nvcc
RM=rm
LINKER=nvcc
TARGET=t

SRCS=$(wildcard $(SRCDIR)/*.cu)
OBJS=$(addprefix $(OBJDIR)/, $(notdir $(SRCS:.cu=.o)))

LIBS=\
	-lpthread\

CFLAGS=\
	-c\
	-O0\
	-g\
	-std=c++17\
	-I$(SRCDIR)\

LFLAGS=\

all: $(OUTDIR)/$(TARGET)

$(OBJDIR)/%.o: $(SRCDIR)/%.cu
	@[ -d $(OBJDIR) ] || mkdir -p $(OBJDIR)
	$(CC) $(CFLAGS) -o $@ -c $<

$(OUTDIR)/$(TARGET): $(OBJS)
	@[ -d $(OUTDIR) ] || mkdir -p $(OUTDIR)
	$(LINKER) $(LFLAGS) $(LIBDIRS) $^ -o $@ $(LIBS)

clean:
	$(RM) $(OBJS) $(OUTDIR)/$(TARGET)

.PHONY: all clean
