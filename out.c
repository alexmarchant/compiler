#include <stdio.h>
#include <stdlib.h>
#include "runtime/runtime.h"

typedef struct _Person {
	String* name;
} Person;

Person* Person__make() {
	Person* val = malloc(sizeof(Person));
	if (!val) {
		printf("Error allocating memory");
		exit(1);
	}
	return val;
}

String* Person__toString(Person* self) {
	return self->name;
}

typedef struct _Point {
	int x;
	int y;
} Point;

Point* Point__make() {
	Point* val = malloc(sizeof(Point));
	if (!val) {
		printf("Error allocating memory");
		exit(1);
	}
	return val;
}

typedef struct _Square {
	int width;
	int height;
} Square;

Square* Square__make() {
	Square* val = malloc(sizeof(Square));
	if (!val) {
		printf("Error allocating memory");
		exit(1);
	}
	return val;
}

int Square__area(Square* self) {
	return self->width * self->height;
}

int main() {
	Person* person = Person__make();
	person->name = String__make("Alex");
	printf("%s %s\n", String__make("name:")->value, Person__toString(person)->value);
	Square* square = Square__make();
	square->width = 5;
	square->height = 5;
	int area = Square__area(square);
	printf("%s %d\n", String__make("area:")->value, area);
	String* msg = String__make("Hello, World!");
	printf("%s\n", msg->value);
	return 0;
}

