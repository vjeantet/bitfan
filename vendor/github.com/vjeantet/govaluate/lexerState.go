package govaluate

type lexerState struct {
	isEOF          bool
	kind           TokenKind
	validNextKinds []TokenKind
}

// lexer states.
// Constant for all purposes except compiler.
var validLexerStates = []lexerState{

	{

		kind:  CLAUSE,
		isEOF: false,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			BOOLEAN,
			VARIABLE,
			STRING,
			TIME,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},

	{

		kind:  CLAUSE_CLOSE,
		isEOF: true,
		validNextKinds: []TokenKind{

			COMPARATOR,
			MODIFIER,
			PREFIX,
			NUMERIC,
			BOOLEAN,
			VARIABLE,
			STRING,
			TIME,
			CLAUSE,
			CLAUSE_CLOSE,
			LOGICALOP,
		},
	},

	{

		kind:  NUMERIC,
		isEOF: true,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
		},
	},
	{

		kind:  BOOLEAN,
		isEOF: true,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
		},
	},
	{

		kind:  PATTERN,
		isEOF: true,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
		},
	},
	{

		kind:  ARRAY,
		isEOF: true,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
		},
	},
	{

		kind:  STRING,
		isEOF: true,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
		},
	},
	{

		kind:  TIME,
		isEOF: true,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
		},
	},
	{

		kind:  VARIABLE,
		isEOF: true,
		validNextKinds: []TokenKind{

			MODIFIER,
			COMPARATOR,
			LOGICALOP,
			CLAUSE_CLOSE,
		},
	},
	{

		kind:  MODIFIER,
		isEOF: false,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			VARIABLE,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},
	{

		kind:  COMPARATOR,
		isEOF: false,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			BOOLEAN,
			VARIABLE,
			STRING,
			TIME,
			CLAUSE,
			CLAUSE_CLOSE,
			PATTERN,
			ARRAY,
		},
	},
	{

		kind:  LOGICALOP,
		isEOF: false,
		validNextKinds: []TokenKind{

			PREFIX,
			NUMERIC,
			BOOLEAN,
			VARIABLE,
			STRING,
			TIME,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},
	{

		kind:  PREFIX,
		isEOF: false,
		validNextKinds: []TokenKind{

			NUMERIC,
			BOOLEAN,
			VARIABLE,
			CLAUSE,
			CLAUSE_CLOSE,
		},
	},
}

func (this lexerState) canTransitionTo(kind TokenKind) bool {

	for _, validKind := range this.validNextKinds {

		if validKind == kind {
			return true
		}
	}

	return false
}
