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
#include <string.h>

static GType * 
alloc_types(int n) {
	return ((GType *)g_new0(GType, n));
}

static void
set_type(GType *types, int n, GType t)
{
	types[n] = t;
}

static GtkTreeViewColumn *
_gtk_tree_view_column_new_with_attributes_one(const gchar *title,
    GtkCellRenderer *renderer, const gchar *attribute, gint column)
{
	GtkTreeViewColumn	*tvc;

	tvc = gtk_tree_view_column_new_with_attributes(title, renderer,
	    attribute, column, NULL);
	return (tvc);
}

static GtkWidget *
_gtk_message_dialog_new(GtkWindow *parent, GtkDialogFlags flags,
    GtkMessageType type, GtkButtonsType buttons, char *msg)
{
	GtkWidget		*w;

	w = gtk_message_dialog_new(parent, flags, type, buttons, "%s", msg);
	return (w);
}

static gchar *
error_get_message(GError *error)
{
	return error->message;
}

static const gchar *
object_get_class_name(GObject *object)
{
	return G_OBJECT_CLASS_NAME(G_OBJECT_GET_CLASS(object));
}
