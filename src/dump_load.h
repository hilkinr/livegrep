/********************************************************************
 * livegrep -- dump_load.h
 * Copyright (c) 2011-2013 Nelson Elhage
 *
 * This program is free software. You may use, redistribute, and/or
 * modify it under the terms listed in the COPYING file.
 ********************************************************************/
#ifndef CODESEARCH_DUMP_LOAD_H
#define CODESEARCH_DUMP_LOAD_H

const uint32_t kIndexMagic   = 0xc0d35eac;
const uint32_t kIndexVersion = 10;

struct index_header {
    uint32_t magic;
    uint32_t version;
    uint32_t chunk_size;

    uint32_t ntrees;
    uint64_t refs_off;

    uint32_t nfiles;
    uint64_t files_off;

    uint32_t nchunks;
    uint64_t chunks_off;

    uint32_t ncontent;
    uint64_t content_off;
} __attribute__((packed));

struct chunk_header {
    uint64_t data_off;
    uint64_t files_off;
    uint32_t size;
    uint32_t nfiles;
} __attribute__((packed));

struct content_chunk_header {
    uint64_t file_off;
    uint32_t size;
} __attribute__((packed));

#endif
