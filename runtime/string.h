#include <stdio.h>

#ifndef STRING_H
#define STRING_H

typedef struct _String {
    char* value;
} String;

String* String__make(char* value);

#endif