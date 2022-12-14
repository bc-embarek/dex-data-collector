pragma solidity ^0.8.0;

import "./interfaces/IERC20.sol";
import "./interfaces/IUniswapV2Pair.sol";
import "./interfaces/IUniswapV2Factory.sol";
import "./interfaces/IUniswapV2Router01.sol";

contract DataCollector {
    struct Pool {
        address addr;
        address token0;
        address token1;
    }

    struct Token {
        address addr;
        string name;
        string symbol;
        uint8 decimal;
    }

    function getPoolAddresses(address _factory, uint start, uint size) external view returns (Pool[] memory) {
        IUniswapV2Factory factory = IUniswapV2Factory(_factory);

        uint pairLength = factory.allPairsLength();
        uint end = pairLength;
        if (pairLength > (start + size)) {
            end = start + size;
        }

        Pool[] memory pools = new Pool[](end - start);
        for (uint idx = start; idx < end; idx++) {
            IUniswapV2Pair pair = IUniswapV2Pair(factory.allPairs(idx));
            pools[idx - start] = Pool(address(pair), pair.token0(), pair.token1());
        }
        return pools;
    }

    function getTokenInfo(address[] memory _tokens) external view returns (Token[] memory) {
        Token[] memory tokens = new Token[](_tokens.length);
        for (uint i = 0; i < _tokens.length; i++) {
            IERC20 token = IERC20(_tokens[i]);
            tokens[i] = Token(_tokens[i], token.name(), token.symbol(), token.decimals());
        }
        return tokens;
    }
}