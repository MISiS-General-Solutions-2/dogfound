/*
This is a Optical-Character-Recognition program
Copyright (C) 2000-2017 Joerg Schulenburg License: GPL2

*/
#ifndef _OCR0_DBG_H
#define _OCR0_DBG_H

//#include "gocr.h"
//#include "unicode_defs.h"

#define MM                                                                 \
  {                                                                        \
    IFV fprintf(stderr, "\nDBG %c L%04d (%d,%d): ", (char)c_ask, __LINE__, \
                box1->x0, box1->y0);                                       \
  }

// new debug mode (0.41) explains why char is declined or accepted as ABC...
//   the output can be filtered by external scripts
//   ToDo: we could reduce output to filter string, +move to ocr0dbg.h
#ifndef DO_DEBUG   /* can be defined outside (configure --with-debug) */
#define DO_DEBUG 0 /* 0 is the default, 1 is for debug+developping */
#endif

#define Setac(box1, ac, ad, job) setac(box1, ac, ad, job)
#define Break break
#define MSG(x)
#define DBG(x)

#endif
