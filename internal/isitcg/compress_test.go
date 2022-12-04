package isitcg

// [Theory]
// [InlineData("hello world", "y0jNyclXKM8vykkBAA")]
// [InlineData(
// 	"{\"n\":\"Test Ingredients\",\"i\":[\"Aqua (Purified Water)\",\"Sodium Cocoamphoacetate\",\"Lauryl Glucoside\",\"Sodium Cocoyl Glutamate\",\"Sodium Lauryl Glucose Carboxylate\",\"Decyl Glucoside\",\"Cocamidopropyl Betaine\",\"Glycerin\",\"Panthenol\",\"Potassium Sorbate\",\"Citric Acid\",\"Polysorbate 20\",\"Phenoxyethanol\",\"Menthol\",\"Mentha Piperita (Peppermint) Oil\",\"Aloe Barbadensis (Aloe Vera) Leaf Juice\",\"Carthamus Tinctorius (Safflower) Oil\",\"Achillea Millefolium Extract\",\"Chamomilla Recutita (Matricaria) Flower Extract\",\"Equisetum Arvense Extract\",\"Eucalyptus Globulus (Eucalyptus) Oil\",\"Lavandula Angustifolia (Lavender) Leaf Extract\",\"Melaleuca Alternifolia (Tea Tree) Leaf Oil\",\"Rosmarinus Officinalis (Rosemary) Leaf Extract\",\"Salvia Officinalis (Sage) Leaf Extract\",\"Thymus Vulgaris (Thyme) Extract\",\"Tussilago Farfara (Coltsfoot) Flower Extract\",\"Urtica Dioica (Nettle) Extract\"]}",
// 	"ZVLLbttADPyVhU420EPRY26K4xgt7MaI3fRQ9ECvKIkAtVT24Voo-u_lSnGjNCcRHHJmONrfhStuiiOGaD67xmNF6GIoPhRU3PwoyucEZrFPnmrCynyHiH6p4EEqSp1ZiRXo-lbAYlRMkS0kP7DZcLISqMqt2fCEROim4RfkzQ6aFfiTXAaeZu7Q_senqtBRJb2XXqFblSaXhTY8WPTktNyDiy064VxLhBCy34P408S6oujJmtJSNU7wECbMfPqYG3n3MmBsYeLYaSrtyDZWYPbUq1TM6WCvZUcuLs0DZcGSBc2tXgEVukDBLMbOE3pYmi1Cbb4kstmxnqoSXQrmSM5G8aTl4gB1zfJLo74S2paYEcwuf2rhfMz6Ej3YmFmUQjqFwDyiTXG0tYN8IXhSzfuRbbaxfk4UMCpL6c_qEedYssBDH9XJhuWUOFta_2tePW3hDK5Kqlm6JoVI2ZbGoX10VfY-XvrqcocMjMpjStZn5K4LR73r6BFfFqYIHyV06t2p9kNdkyUHnIPUPiowvGM_AJ9V_s3wAZor66uNYzvkvJ8SNyqglLmhY7OJpK-FoRFzD74Gr0ethGOoRfQXv8vym4-as7kjyZ_FV4yRZ3w___wF"
// )]
// public void CanCompress(string input, string expected)
// {
// 	var actual = StringCompression.Compress(input);
// 	Assert.Equal(expected, actual);
// }

// [Theory]
// [InlineData("hello world", "y0jNyclXKM8vykkBAA")]
// [InlineData(
// 	"{\"n\":\"Test Ingredients\",\"i\":[\"Aqua (Purified Water)\",\"Sodium Cocoamphoacetate\",\"Lauryl Glucoside\",\"Sodium Cocoyl Glutamate\",\"Sodium Lauryl Glucose Carboxylate\",\"Decyl Glucoside\",\"Cocamidopropyl Betaine\",\"Glycerin\",\"Panthenol\",\"Potassium Sorbate\",\"Citric Acid\",\"Polysorbate 20\",\"Phenoxyethanol\",\"Menthol\",\"Mentha Piperita (Peppermint) Oil\",\"Aloe Barbadensis (Aloe Vera) Leaf Juice\",\"Carthamus Tinctorius (Safflower) Oil\",\"Achillea Millefolium Extract\",\"Chamomilla Recutita (Matricaria) Flower Extract\",\"Equisetum Arvense Extract\",\"Eucalyptus Globulus (Eucalyptus) Oil\",\"Lavandula Angustifolia (Lavender) Leaf Extract\",\"Melaleuca Alternifolia (Tea Tree) Leaf Oil\",\"Rosmarinus Officinalis (Rosemary) Leaf Extract\",\"Salvia Officinalis (Sage) Leaf Extract\",\"Thymus Vulgaris (Thyme) Extract\",\"Tussilago Farfara (Coltsfoot) Flower Extract\",\"Urtica Dioica (Nettle) Extract\"]}",
// 	"ZVLLbttADPyVhU420EPRY26K4xgt7MaI3fRQ9ECvKIkAtVT24Voo-u_lSnGjNCcRHHJmONrfhStuiiOGaD67xmNF6GIoPhRU3PwoyucEZrFPnmrCynyHiH6p4EEqSp1ZiRXo-lbAYlRMkS0kP7DZcLISqMqt2fCEROim4RfkzQ6aFfiTXAaeZu7Q_senqtBRJb2XXqFblSaXhTY8WPTktNyDiy064VxLhBCy34P408S6oujJmtJSNU7wECbMfPqYG3n3MmBsYeLYaSrtyDZWYPbUq1TM6WCvZUcuLs0DZcGSBc2tXgEVukDBLMbOE3pYmi1Cbb4kstmxnqoSXQrmSM5G8aTl4gB1zfJLo74S2paYEcwuf2rhfMz6Ej3YmFmUQjqFwDyiTXG0tYN8IXhSzfuRbbaxfk4UMCpL6c_qEedYssBDH9XJhuWUOFta_2tePW3hDK5Kqlm6JoVI2ZbGoX10VfY-XvrqcocMjMpjStZn5K4LR73r6BFfFqYIHyV06t2p9kNdkyUHnIPUPiowvGM_AJ9V_s3wAZor66uNYzvkvJ8SNyqglLmhY7OJpK-FoRFzD74Gr0ethGOoRfQXv8vym4-as7kjyZ_FV4yRZ3w___wF"
// )]
// public void CanDecompress(string expected, string input)
// {
// 	var actual = StringCompression.Decompress(input);
// 	Assert.Equal(expected, actual);
// }
