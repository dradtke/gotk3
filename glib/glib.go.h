/*
 * Copyright (c) 2013 Conformal Systems <info@conformal.com>
 *
 * This file originated from: http://opensource.conformal.com/
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>

static GType
_g_type_from_instance(gpointer instance)
{
	return (G_TYPE_FROM_INSTANCE(instance));
}

/* Wrapper to avoid variable arg list */
static void
_g_object_set_one(gpointer object, const gchar *property_name, void *val)
{
	g_object_set(object, property_name, *(gpointer **)val, NULL);
}

static GValue *
alloc_gvalue_list(int n)
{
	GValue		*valv;

	valv = g_new0(GValue, n);
	return (valv);
}

static void
val_list_insert(GValue *valv, int i, GValue *val)
{
	valv[i] = *val;
}

/*
 * GValue
 */

static GValue *
_g_value_alloc()
{
	return (g_new0(GValue, 1));
}

static GValue *
_g_value_init(GType g_type)
{
	GValue          *value;

	value = g_new0(GValue, 1);
	return (g_value_init(value, g_type));
}

static gboolean
_g_is_value(GValue *val)
{
	return (G_IS_VALUE(val));
}

static GType
_g_value_type(GValue *val)
{
	return (G_VALUE_TYPE(val));
}

static GType
_g_value_fundamental(GType type)
{
	return (G_TYPE_FUNDAMENTAL(type));
}

/*
 * Closures
 */

extern void goMarshal(GClosure *closure, GValue *return_value, guint n_param_values, GValue *param_values, gpointer invocation_hint, gpointer marshal_data);

static GClosure *
_g_closure_new()
{
	GClosure *closure = g_closure_new_simple(sizeof(GClosure), NULL);
	g_closure_set_marshal(closure, (GClosureMarshal)(goMarshal));
	return closure;
}

/*
 * Variant types
 */

static const GVariantType *
_g_variant_type_boolean()
{
	return (G_VARIANT_TYPE_BOOLEAN);
}

static const GVariantType *
_g_variant_type_byte()
{
	return (G_VARIANT_TYPE_BYTE);
}

static const GVariantType *
_g_variant_type_int16()
{
	return (G_VARIANT_TYPE_INT16);
}

static const GVariantType *
_g_variant_type_uint16()
{
	return (G_VARIANT_TYPE_UINT16);
}

static const GVariantType *
_g_variant_type_int32()
{
	return (G_VARIANT_TYPE_INT32);
}

static const GVariantType *
_g_variant_type_uint32()
{
	return (G_VARIANT_TYPE_UINT32);
}

static const GVariantType *
_g_variant_type_int64()
{
	return (G_VARIANT_TYPE_INT64);
}

static const GVariantType *
_g_variant_type_uint64()
{
	return (G_VARIANT_TYPE_UINT64);
}

static const GVariantType *
_g_variant_type_handle()
{
	return (G_VARIANT_TYPE_HANDLE);
}

static const GVariantType *
_g_variant_type_double()
{
	return (G_VARIANT_TYPE_DOUBLE);
}

static const GVariantType *
_g_variant_type_string()
{
	return (G_VARIANT_TYPE_STRING);
}

static const GVariantType *
_g_variant_type_object_path()
{
	return (G_VARIANT_TYPE_OBJECT_PATH);
}

static const GVariantType *
_g_variant_type_signature()
{
	return (G_VARIANT_TYPE_SIGNATURE);
}

static const GVariantType *
_g_variant_type_variant()
{
	return (G_VARIANT_TYPE_VARIANT);
}

static const GVariantType *
_g_variant_type_any()
{
	return (G_VARIANT_TYPE_ANY);
}

static const GVariantType *
_g_variant_type_basic()
{
	return (G_VARIANT_TYPE_BASIC);
}

static const GVariantType *
_g_variant_type_maybe()
{
	return (G_VARIANT_TYPE_MAYBE);
}

static const GVariantType *
_g_variant_type_array()
{
	return (G_VARIANT_TYPE_ARRAY);
}

static const GVariantType *
_g_variant_type_tuple()
{
	return (G_VARIANT_TYPE_TUPLE);
}

static const GVariantType *
_g_variant_type_unit()
{
	return (G_VARIANT_TYPE_UNIT);
}

static const GVariantType *
_g_variant_type_dict_entry()
{
	return (G_VARIANT_TYPE_DICT_ENTRY);
}

static const GVariantType *
_g_variant_type_dictionary()
{
	return (G_VARIANT_TYPE_DICTIONARY);
}

static const GVariantType *
_g_variant_type_string_array()
{
	return (G_VARIANT_TYPE_STRING_ARRAY);
}

static const GVariantType *
_g_variant_type_object_path_array()
{
	return (G_VARIANT_TYPE_OBJECT_PATH_ARRAY);
}

static const GVariantType *
_g_variant_type_bytestring()
{
	return (G_VARIANT_TYPE_BYTESTRING);
}

static const GVariantType *
_g_variant_type_bytestring_array()
{
	return (G_VARIANT_TYPE_BYTESTRING_ARRAY);
}

static const GVariantType *
_g_variant_type_vardict()
{
	return (G_VARIANT_TYPE_VARDICT);
}
