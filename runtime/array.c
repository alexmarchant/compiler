#include <stdio.h>
#include <stdlib.h>
#include <strings.h>

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

CharArray* char_array_make() {
  CharArray* array = malloc(sizeof(CharArray));
  if (!array) {
    printf("Error allocating memory");
    exit(1);
  }
  array->elements = malloc(1 * sizeof(char));
  array->length = 0;
  array->capacity = 1;
  return array;
}

IntArray* int_array_make() {
  IntArray* array = malloc(sizeof(IntArray));
  if (!array) {
    printf("Error allocating memory");
    exit(1);
  }
  array->elements = malloc(1 * sizeof(int));
  array->length = 0;
  array->capacity = 1;
  return array;
}

void char_array_push(CharArray* array, char element) {
  // Resize array if needed
  if (array->length == array->capacity) {
    array->capacity = array->capacity * 2; 
    void* new_elements = realloc(
      array->elements,
      array->capacity * sizeof(char)
    );
    if (new_elements) {
      array->elements = new_elements;
    } else {
      printf("Error resizing array");
      exit(1);
    }
  }

  // Push element
  array->elements[array->length] = element;
  array->length++;
}

void int_array_push(IntArray* array, int element) {
  // Resize array if needed
  if (array->length == array->capacity) {
    array->capacity = array->capacity * 2; 
    void* new_elements = realloc(
      array->elements,
      array->capacity * sizeof(int)
    );
    if (new_elements) {
      array->elements = new_elements;
    } else {
      printf("Error resizing array");
      exit(1);
    }
  }

  // Push element
  array->elements[array->length] = element;
  array->length++;
}

void char_array_add(CharArray* array, char* elements) {
  size_t len = strlen(elements);
  for (int i = 0; i < len; i++) {
    char_array_push(array, elements[i]);
  }
}

void int_array_add(IntArray* array, int* elements) {
  size_t len = sizeof(elements) / sizeof(int);
  for (int i = 0; i < len; i++) {
    int_array_push(array, elements[i]);
  }
}

void char_array_print(CharArray* array) {
  printf("\"%s\"\n", array->elements);
}

void int_array_print(IntArray* array) {
  printf("[ ");
  for (int i = 0; i < array->length; i++) {
    printf("%d", ((int*) array->elements)[i]);
    if (i != array->length - 1) {
      printf(", ");
    }
  }
  printf(" ]\n");
}