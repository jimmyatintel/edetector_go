#include "warp_AES.hxx"
#include "PyAES.h"
// g++ -shared -o mylib.so -fPIC warp_AES.cxx
const char *key = "AES Encrypt Decrypt";

void DecryptBuffer_cplus(const char* data,int size, char* outdata) {
  CAES *myAES = new CAES();
  myAES->SetKeys(KeySize::BIT128,key);
  std::vector<char> INPUT(data, data+(size));
  std::vector<char> OUTPUT = myAES->DecryptBuffer(INPUT);
  for (int i = 0; i < OUTPUT.size(); i++) {
    outdata[i] = OUTPUT[i];
  }
  delete myAES;
}

void EncryptBuffer_cplus(const char* data,int size, char* outdata) {
  CAES *myAES = new CAES();
  myAES->SetKeys(KeySize::BIT128,key);
  std::vector<char> INPUT(data, data+(size));
  std::vector<char> OUTPUT = myAES->EncryptBuffer(INPUT);
  for (int i = 0; i < OUTPUT.size(); i++) {
    outdata[i] = OUTPUT[i];
  }
  delete myAES;
}