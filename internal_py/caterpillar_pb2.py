# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: caterpillar.proto

from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='caterpillar.proto',
  package='caterpillar',
  syntax='proto3',
  serialized_options=b'Z*github.com/wpwilson10/caterpillar/protobuf',
  serialized_pb=b'\n\x11\x63\x61terpillar.proto\x12\x0b\x63\x61terpillar\" \n\x10NewspaperRequest\x12\x0c\n\x04link\x18\x01 \x01(\t\"p\n\x0eNewspaperReply\x12\x0c\n\x04link\x18\x01 \x01(\t\x12\r\n\x05title\x18\x02 \x01(\t\x12\x0c\n\x04text\x18\x03 \x01(\t\x12\x11\n\tcanonical\x18\x04 \x01(\t\x12\x0f\n\x07pubdate\x18\x05 \x01(\t\x12\x0f\n\x07\x61uthors\x18\x06 \x03(\t\"\x1b\n\x0bTextRequest\x12\x0c\n\x04text\x18\x01 \x01(\t\"\"\n\rSentenceReply\x12\x11\n\tsentences\x18\x01 \x03(\t\"1\n\x0cSummaryReply\x12\x0f\n\x07summary\x18\x01 \x01(\t\x12\x10\n\x08keywords\x18\x02 \x03(\t2\xdf\x01\n\x0b\x43\x61terpillar\x12I\n\tNewspaper\x12\x1d.caterpillar.NewspaperRequest\x1a\x1b.caterpillar.NewspaperReply\"\x00\x12\x43\n\tSentences\x12\x18.caterpillar.TextRequest\x1a\x1a.caterpillar.SentenceReply\"\x00\x12@\n\x07Summary\x12\x18.caterpillar.TextRequest\x1a\x19.caterpillar.SummaryReply\"\x00\x42,Z*github.com/wpwilson10/caterpillar/protobufb\x06proto3'
)




_NEWSPAPERREQUEST = _descriptor.Descriptor(
  name='NewspaperRequest',
  full_name='caterpillar.NewspaperRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='link', full_name='caterpillar.NewspaperRequest.link', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=34,
  serialized_end=66,
)


_NEWSPAPERREPLY = _descriptor.Descriptor(
  name='NewspaperReply',
  full_name='caterpillar.NewspaperReply',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='link', full_name='caterpillar.NewspaperReply.link', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='title', full_name='caterpillar.NewspaperReply.title', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='text', full_name='caterpillar.NewspaperReply.text', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='canonical', full_name='caterpillar.NewspaperReply.canonical', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='pubdate', full_name='caterpillar.NewspaperReply.pubdate', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='authors', full_name='caterpillar.NewspaperReply.authors', index=5,
      number=6, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=68,
  serialized_end=180,
)


_TEXTREQUEST = _descriptor.Descriptor(
  name='TextRequest',
  full_name='caterpillar.TextRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='text', full_name='caterpillar.TextRequest.text', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=182,
  serialized_end=209,
)


_SENTENCEREPLY = _descriptor.Descriptor(
  name='SentenceReply',
  full_name='caterpillar.SentenceReply',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='sentences', full_name='caterpillar.SentenceReply.sentences', index=0,
      number=1, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=211,
  serialized_end=245,
)


_SUMMARYREPLY = _descriptor.Descriptor(
  name='SummaryReply',
  full_name='caterpillar.SummaryReply',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='summary', full_name='caterpillar.SummaryReply.summary', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='keywords', full_name='caterpillar.SummaryReply.keywords', index=1,
      number=2, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=247,
  serialized_end=296,
)

DESCRIPTOR.message_types_by_name['NewspaperRequest'] = _NEWSPAPERREQUEST
DESCRIPTOR.message_types_by_name['NewspaperReply'] = _NEWSPAPERREPLY
DESCRIPTOR.message_types_by_name['TextRequest'] = _TEXTREQUEST
DESCRIPTOR.message_types_by_name['SentenceReply'] = _SENTENCEREPLY
DESCRIPTOR.message_types_by_name['SummaryReply'] = _SUMMARYREPLY
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

NewspaperRequest = _reflection.GeneratedProtocolMessageType('NewspaperRequest', (_message.Message,), {
  'DESCRIPTOR' : _NEWSPAPERREQUEST,
  '__module__' : 'caterpillar_pb2'
  # @@protoc_insertion_point(class_scope:caterpillar.NewspaperRequest)
  })
_sym_db.RegisterMessage(NewspaperRequest)

NewspaperReply = _reflection.GeneratedProtocolMessageType('NewspaperReply', (_message.Message,), {
  'DESCRIPTOR' : _NEWSPAPERREPLY,
  '__module__' : 'caterpillar_pb2'
  # @@protoc_insertion_point(class_scope:caterpillar.NewspaperReply)
  })
_sym_db.RegisterMessage(NewspaperReply)

TextRequest = _reflection.GeneratedProtocolMessageType('TextRequest', (_message.Message,), {
  'DESCRIPTOR' : _TEXTREQUEST,
  '__module__' : 'caterpillar_pb2'
  # @@protoc_insertion_point(class_scope:caterpillar.TextRequest)
  })
_sym_db.RegisterMessage(TextRequest)

SentenceReply = _reflection.GeneratedProtocolMessageType('SentenceReply', (_message.Message,), {
  'DESCRIPTOR' : _SENTENCEREPLY,
  '__module__' : 'caterpillar_pb2'
  # @@protoc_insertion_point(class_scope:caterpillar.SentenceReply)
  })
_sym_db.RegisterMessage(SentenceReply)

SummaryReply = _reflection.GeneratedProtocolMessageType('SummaryReply', (_message.Message,), {
  'DESCRIPTOR' : _SUMMARYREPLY,
  '__module__' : 'caterpillar_pb2'
  # @@protoc_insertion_point(class_scope:caterpillar.SummaryReply)
  })
_sym_db.RegisterMessage(SummaryReply)


DESCRIPTOR._options = None

_CATERPILLAR = _descriptor.ServiceDescriptor(
  name='Caterpillar',
  full_name='caterpillar.Caterpillar',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=299,
  serialized_end=522,
  methods=[
  _descriptor.MethodDescriptor(
    name='Newspaper',
    full_name='caterpillar.Caterpillar.Newspaper',
    index=0,
    containing_service=None,
    input_type=_NEWSPAPERREQUEST,
    output_type=_NEWSPAPERREPLY,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='Sentences',
    full_name='caterpillar.Caterpillar.Sentences',
    index=1,
    containing_service=None,
    input_type=_TEXTREQUEST,
    output_type=_SENTENCEREPLY,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='Summary',
    full_name='caterpillar.Caterpillar.Summary',
    index=2,
    containing_service=None,
    input_type=_TEXTREQUEST,
    output_type=_SUMMARYREPLY,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_CATERPILLAR)

DESCRIPTOR.services_by_name['Caterpillar'] = _CATERPILLAR

# @@protoc_insertion_point(module_scope)
