#include <cmath>
#include <stdint.h>
#include <string>
#include <algorithm>
#include "cryptonote_basic/cryptonote_basic.h"
#include "cryptonote_basic/cryptonote_format_utils.h"
#include "cryptonote_basic/blobdatatype.h"
#include "crypto/crypto.h"
#include "crypto/hash.h"
#include "common/base58.h"
#include "serialization/binary_utils.h"

using namespace cryptonote;

#define HASH_SIZE 32

/*******************************************************************************************/
/* Helper functions for merged mining - for merkle tree branch compuation. */

static size_t tree_depth(size_t count)
{
  size_t i;
  size_t depth = 0;
  assert(count > 0);
  for (i = sizeof(size_t) << 2; i > 0; i >>= 1)
  {
    if (count >> i > 0)
    {
      count >>= i;
      depth += i;
    }
  }
  return depth;
}

static void tree_branch(const char (*hashes)[HASH_SIZE], size_t count, char (*branch)[HASH_SIZE])
{
  size_t i, j;
  size_t cnt = 1;
  size_t depth = 0;
  char (*ints)[HASH_SIZE];
  assert(count > 0);
  for (i = sizeof(size_t) << 2; i > 0; i >>= 1)
  {
    if (cnt << i <= count)
    {
      cnt <<= i;
      depth += i;
    }
  }
  assert(cnt == 1ULL << depth);
  assert(depth == tree_depth(count));
  ints = reinterpret_cast<char(*)[HASH_SIZE]>(alloca((cnt - 1) * HASH_SIZE));
  memcpy(ints, hashes + 1, (2 * cnt - count - 1) * HASH_SIZE);
  for (i = 2 * cnt - count, j = 2 * cnt - count - 1; j < cnt - 1; i += 2, ++j)
  {
    crypto::cn_fast_hash(hashes[i], 2 * HASH_SIZE, ints[j]);
  }
  assert(i == count);
  while (depth > 0)
  {
    assert(cnt == 1ULL << depth);
    cnt >>= 1;
    --depth;
    memcpy(branch[depth], ints[0], HASH_SIZE);
    for (i = 1, j = 0; j < cnt - 1; i += 2, ++j)
    {
      crypto::cn_fast_hash(ints[i], 2 * HASH_SIZE, ints[j]);
    }
  }
}

static void tree_hash_from_branch(const char (*branch)[HASH_SIZE], size_t depth, const char* leaf, const void* path, char* root_hash)
{
  if (depth == 0)
  {
    memcpy(root_hash, leaf, HASH_SIZE);
  }
  else
  {
    char buffer[2][HASH_SIZE];
    int from_leaf = 1;
    char *leaf_path, *branch_path;
    while (depth > 0)
    {
      --depth;
      if (path && (((const char*) path)[depth >> 3] & (1 << (depth & 7))) != 0)
      {
        leaf_path = buffer[1];
        branch_path = buffer[0];
      }
      else
      {
        leaf_path = buffer[0];
        branch_path = buffer[1];
      }
      if (from_leaf)
      {
        memcpy(leaf_path, leaf, HASH_SIZE);
        from_leaf = 0;
      }
      else
      {
        crypto::cn_fast_hash(buffer, 2 * HASH_SIZE, leaf_path);
      }
      memcpy(branch_path, branch[depth], HASH_SIZE);
    }
    crypto::cn_fast_hash(buffer, 2 * HASH_SIZE, root_hash);
  }
}


#if 0
class CAuxPow
{
public:
    // Cryptnote coinbase tx, which contains the block hash of kevacoin block.
    cryptonote::transaction miner_tx;

    // Merkle branch is used to establish that miner_tx is part of the
    // merkel tree whose root is merkle_root.
    std::vector<crypto::hash>  merkle_branch;

    // load
    template <template <bool> class Archive>
    bool do_serialize(Archive<false>& ar)
    {
        FIELD(miner_tx)
        FIELD(merkle_branch)
        return true;
    }

    // store
    template <template <bool> class Archive>
    bool do_serialize(Archive<true>& ar)
    {
        FIELD(miner_tx)
        FIELD(merkle_branch)
        return true;
    }
};
#endif

class CAuxPow
{
public:
    // Cryptnote coinbase tx, which contains the block hash of kevacoin block.
    cryptonote::transaction miner_tx;

    // Merkle branch is used to establish that miner_tx is part of the
    // merkel tree whose root is merkle_root.
    std::vector<crypto::hash>  merkle_branch;

    // load
    template <template <bool> class Archive>
    bool do_serialize(Archive<false>& ar)
    {
        FIELD(miner_tx)
        FIELD(merkle_branch)
        return true;
    }

    // store
    template <template <bool> class Archive>
    bool do_serialize(Archive<true>& ar)
    {
        FIELD(miner_tx)
        FIELD(merkle_branch)
        return true;
    }
};


// Merkle branch and the whole miner tx.
static blobdata get_block_auxpow_blob(const block& b)
{
    CAuxPow auxPow;
    auxPow.miner_tx = b.miner_tx;

    std::vector<crypto::hash> transactionHashes;
    transactionHashes.push_back(cryptonote::get_transaction_hash(b.miner_tx));
    std::copy(b.tx_hashes.begin(), b.tx_hashes.end(), std::back_inserter(transactionHashes));
    auxPow.merkle_branch.resize(tree_depth(b.tx_hashes.size() + 1));
    tree_branch(reinterpret_cast<const char(*)[HASH_SIZE]>(transactionHashes.data()), transactionHashes.size(),
            reinterpret_cast<char(*)[HASH_SIZE]>(auxPow.merkle_branch.data()));

    return t_serializable_object_to_blob(auxPow);
}

extern "C" uint32_t convert_blob_to_auxpow_blob(const char *blob, uint32_t len, char** out) {
    std::string input = std::string(blob, len);
    std::string output = "";

    block b = AUTO_VAL_INIT(b);
    if (!parse_and_validate_block_from_blob(input, b)) {
        return 0;
    }

    output = get_block_auxpow_blob(b);
    *out = (char*)malloc(output.length());
    output.copy(*out, output.length(), 0);
    return output.length();
}

extern "C" uint32_t convert_blob(const char *blob, uint32_t len, char *out) {
    std::string input = std::string(blob, len);
    blobdata output = "";

    block b = AUTO_VAL_INIT(b);
    if (!cryptonote::parse_and_validate_block_from_blob(input, b)) {
        return 0;
    }

    output = cryptonote::get_block_hashing_blob(b);
    output.copy(out, output.length(), 0);
    return output.length();
}

extern "C" bool validate_address(const char *addr, size_t len) {
    std::string input = std::string(addr, len);
    std::string output = "";
    uint64_t prefix;
    return tools::base58::decode_addr(addr, prefix, output);
}


extern "C" void cryptonight_hash(const char* input, char* output, uint32_t len) {
    const int variant = input[0] >= 7 ? input[0] - 6 : 0;
    crypto::cn_slow_hash(input, len, output, variant, 0, 0);
}

extern "C" void cryptonight_fast_hash(const char* input, char* output, uint32_t len) {
    crypto::cn_fast_hash(input, len, output);
}
