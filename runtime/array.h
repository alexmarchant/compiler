#include <stdio.h>

#ifndef ARRAY_H
#define ARRAY_H

typedef struct _CharArray {
  char* elements;
  size_t length;
  size_t capacity;
} CharArray;

typedef struct _IntArray {
  int* elements;
  size_t length;
  size_t capacity;
} IntArray;

CharArray* char_array_make();
IntArray* int_array_make();
void char_array_push(CharArray* array, char* element);
void int_array_push(IntArray* array, int element);
void char_array_add(CharArray* array, char* elements);
void int_array_add(IntArray* array, int* elements);
void char_array_print(CharArray* array);
void int_array_print(IntArray* array);

#endif