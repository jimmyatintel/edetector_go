#ifndef WRAP_AES_H
#define WRAP_AES_H

// __cplusplus gets defined when a C++ compiler processes the file
#ifdef __cplusplus
// extern "C" is needed so the C++ compiler exports the symbols without name
// manging.
extern "C" {
#endif

void DecryptBuffer_cplus(const char* data,int size,char* outdata);

void EncryptBuffer_cplus(const char* data,int size,char* outdata);

#ifdef __cplusplus
}
#endif


#endif
