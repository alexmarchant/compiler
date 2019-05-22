#include <stdio.h>
#include <stdlib.h>
#include <strings.h>

typedef struct _String {
    char* value;
} String;

String* String__make(char* value) {
    String* val = malloc(sizeof(String));
    if (!val) {
        printf("Error allocating memory");
        exit(1);
    }
    val->value = value;
    return val;
}