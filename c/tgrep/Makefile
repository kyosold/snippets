CC		= gcc
CFLAGS	= -g
INCS	= 
LDS		=
SHLD	= $(CC) $(CFLAGS)

OBJ		= tgrep

all: $(OBJ)

$(OBJ): tgrep.o utils.o xmalloc.o
	$(SHLD) -o $(OBJ) tgrep.o utils.o xmalloc.o

%.o: %.c
	$(SHLD) -c $<

clean:
	rm -f *.o $(OBJ)
