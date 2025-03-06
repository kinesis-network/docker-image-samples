OUTDIR=bin
OBJDIR=obj
SRCDIR=src

CC=g++
RM=rm
LINKER=g++
TARGET=t

SRCS=$(wildcard $(SRCDIR)/*.cpp)
OBJS=$(addprefix $(OBJDIR)/, $(notdir $(SRCS:.cpp=.o)))

LIBS=\
	-lpthread\

CFLAGS=\
	-c\
	-O0\
	-Wall\
	-g\
	-std=c++17\

LFLAGS=\

all: $(OUTDIR)/$(TARGET)

$(OBJDIR)/%.o: $(SRCDIR)/%.cpp
	@[ -d $(OBJDIR) ] || mkdir -p $(OBJDIR)
	$(CC) $(CFLAGS) -o $@ -c $<

$(OUTDIR)/$(TARGET): $(OBJS)
	@[ -d $(OUTDIR) ] || mkdir -p $(OUTDIR)
	$(LINKER) $(LFLAGS) $(LIBDIRS) $^ -o $@ $(LIBS)

clean:
	$(RM) $(OBJS) $(OUTDIR)/$(TARGET)

.PHONY: all clean
