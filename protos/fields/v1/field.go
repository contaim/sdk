// Copyright (c) 2021 Contaim, LLC
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package fields

func NewFields(fields ...*Field) *Spec {
	return &Spec{
		Field: fields,
	}
}

func NewTextField(name, label string) *Field {
	return &Field{
		Field: &Field_Text{
			Text: &Text{
				Base: &Base{
					Name:  name,
					Label: label,
				},
			},
		},
	}
}
