> Content sourced from [Algebrica](https://algebrica.org/cosine-similarity/) — CC BY-NC 4.0

## What is cosine similarity

When working with large collections of text, it becomes necessary to have tools that allow a computer to evaluate how similar two documents are to each other. One of the most widely used approaches in this context is cosine similarity, a measure rooted in linear algebra that quantifies the relationship between two vectors in terms of the angle they form.

The central idea is that any document can be represented as a [vector](<../vectors/>) in a high-dimensional space, where each component corresponds to a feature of the text, such as the frequency of a given word or its contextual relevance. Once two documents have been encoded in this way, their degree of similarity can be assessed by examining the geometric relationship between the corresponding vectors. What distinguishes cosine similarity from simpler distance measures is that it focuses entirely on the orientation of the vectors, disregarding their magnitude: two documents will be considered similar if they point in roughly the same direction in the feature space, regardless of how long they are.

The name of the measure comes directly from the [cosine function](<../cosine-function/>) : cosine similarity is defined as the cosine of the angle between the two vectors, and its value is expressed by the following formula, which is discussed in detail in the next section:

\\[C_s(V_x, V_y) = \frac{V_x \cdot V_y}{\|V_x\| \, \|V_y\|} \\]

Although more sophisticated methods exist for measuring semantic similarity, such as neural embeddings or transformer-based language models, cosine similarity offers a compelling balance between mathematical simplicity and practical effectiveness. It is widely used in recommendation systems, automatic text classification, and semantic search, where the ability to quickly assess the relationship between documents is essential. In what follows, we will see how to compute cosine similarity step by step, starting from a concrete example involving a small set of sentences.

> It is worth keeping in mind that cosine similarity carries no semantic understanding of language whatsoever. The measure does not grasp the actual meaning of words or sentences but operates purely on a geometric notion of proximity between vectors. Two documents may be judged similar because their numerical representations point in the same direction even when their content is conceptually unrelated, and conversely, documents with closely related meanings may appear distant if their vector representations differ in orientation.

## Distance between vectors and its limitations

Before introducing cosine similarity, it is worth examining a more elementary approach to measuring how close two vectors are, namely the Euclidean distance. Given two vectors \\(\mathbf{u} = (u_1, u_2, \ldots, u_n)\\) and \\(\mathbf{v} = (v_1, v_2, \ldots, v_n)\\), the Euclidean distance between them is defined as follows:

\\[d(\mathbf{u}, \mathbf{v}) = \sqrt{\sum_{i=1}^{n} (u_i - v_i)^2} \\]

This formula measures the straight-line distance between the two points in \\(n\\)-dimensional space that the vectors identify.

![](/diagrams/algebrica/euclidean-distance.png)

A small value of \\(d(\mathbf{u}, \mathbf{v})\\) indicates that the two vectors are geometrically close to each other, while a large value indicates that they are far apart. Since the Euclidean distance is an unbounded quantity, it is sometimes convenient to convert it into a similarity score that takes values in a bounded interval. One common way to do this is through the following transformation:

\\[\text{sim}(\mathbf{u}, \mathbf{v}) = \frac{1}{1 + d(\mathbf{u}, \mathbf{v})} \\]

This expression maps the distance to a value in the interval \\((0, 1]\\). When the two vectors are identical, the distance is zero and the similarity equals \\(1\\). As the distance increases, the similarity decreases monotonically towards \\(0\\), without ever reaching it.


The Euclidean distance, however, has a significant limitation when applied to text analysis: it is sensitive to the magnitude of the vectors, not only to their direction. Two documents discussing exactly the same topics will produce vectors that point in the same direction, but one of them might have a much larger norm simply because it is longer. The Euclidean distance would then indicate that the two documents are far apart, even though their content is essentially identical. Cosine similarity addresses this limitation by normalizing the vectors before comparing them, so that only the angle between the two directions is taken into account.

## How to calculate cosine similarity between two vectors

Cosine similarity is defined as the cosine of the angle formed by two vectors in a given [vector space](<../vector-spaces/>). Given two vectors \\(\mathbf{V_x}\\) and \\(\mathbf{V_y}\\) in \\(\mathbb{R}^n\\), the cosine similarity between them is expressed by the following formula:

\\[C_s(V_x, V_y) = \frac{\displaystyle\sum_{i=1}^{n} V_{x_i} \cdot V_{y_i}}{\sqrt{\displaystyle\sum_{i=1}^{n} (V_{x_i})^2} \cdot \sqrt{\displaystyle\sum_{i=1}^{n} (V_{y_i})^2}} \\]

In compact notation, using the dot product and the Euclidean norm, the same expression takes the following form:

\\[C_s(V_x, V_y) = \frac{V_x \cdot V_y}{\|V_x\| \, \|V_y\|} \\]

In these expressions, \\(V_x \cdot V_y\\) denotes the dot product of the two vectors, while \\(\|V_x\|\\) and \\(\|V_y\|\\) are their respective Euclidean norms. Dividing by the product of the norms is precisely the normalization step that removes the effect of vector magnitude and retains only directional information.

![](/diagrams/algebrica/cosine-similarity.png)

The value of cosine similarity ranges between \\(-1\\) and \\(1\\). In text analysis, where vector components are non-negative by construction, the range is restricted to \\([0, 1]\\). A value close to \\(1\\) indicates that the angle between the two vectors is small, meaning the vectors are nearly parallel and the corresponding documents are highly similar. A value close to \\(0\\) indicates that the vectors are nearly orthogonal, and therefore that the two documents share little to no common content.

> From a purely mathematical standpoint, a value of \\(-1\\) indicates that the two vectors point in exactly opposite directions, forming an angle of \\(180°\\). This case does not arise in text analysis, where all vector components are non-negative, but it remains part of the general mathematical definition.

## Example 1

Consider three sentences for which we wish to determine the degree of mutual similarity:

  * \\(x\\) = I am fond of reading thriller novels.
  * \\(y\\) = I prefer reading thriller novels.
  * \\(z\\) = Yesterday, I arrived late.


A preliminary inspection suggests that sentences \\(x\\) and \\(y\\) share a common theme, while sentence \\(z\\) is clearly unrelated to the other two.

The first step consists in transforming each sentence into a vector by extracting all distinct words and recording their frequency of occurrence. Before doing so, we remove words that carry little semantic content, such as the conjunction “of”, the pronoun “I”, and the verb “to be”. This filtering step, known as stop-word removal, is standard practice in text analysis, particularly when working with large corpora, as it ensures that the resulting vectors reflect only the most meaningful lexical elements.

| arrived| fond| late| novels| prefer| reading| thriller| yesterday  
---|---|---|---|---|---|---|---|---  
\\(V_x\\)| 0| 1| 0| 1| 0| 1| 1| 0  
\\(V_y\\)| 0| 0| 0| 1| 1| 1| 1| 0  
\\(V_z\\)| 1| 0| 1| 0| 0| 0| 0| 1  
  
The result of this process is the previous term-document matrix, where each entry indicates whether a given word is present in the corresponding sentence. The values are binary in this case because each meaningful word appears at most once in each sentence; in general, the entries can also represent raw word frequencies or weighted values such as TF-IDF scores. The vector representation of the three sentences is therefore the following:

  * \\(V_x = [0, 1, 0, 1, 0, 1, 1, 0]\\)
  * \\(V_y = [0, 0, 0, 1, 1, 1, 1, 0]\\)
  * \\(V_z = [1, 0, 1, 0, 0, 0, 0, 1]\\)


We now apply the cosine similarity formula to the pair \\((V_x, V_y)\\), which we expect to yield a high similarity value. The formula to be evaluated is the following:

\\[C_s(V_x, V_y) = \frac{V_x \cdot V_y}{\|V_x\| \, \|V_y\|} \\]

We begin by computing the dot product \\(V_x \cdot V_y\\). Multiplying the corresponding components of the two vectors and summing the results, we obtain:

\\[V_x \cdot V_y = (0 \cdot 0) + (1 \cdot 0) + (0 \cdot 0) + (1 \cdot 1) + (0 \cdot 1) + (1 \cdot 1) + (1 \cdot 1) + (0 \cdot 0) = 3 \\]

We then compute the Euclidean norm of each vector, which appears in the denominator of the formula. The norm of \\(V_x\\) is given by:

\\[\|V_x\| = \sqrt{0^2 + 1^2 + 0^2 + 1^2 + 0^2 + 1^2 + 1^2 + 0^2} = \sqrt{4} = 2 \\]

The norm of \\(V_y\\) is computed in the same way:

\\[\|V_y\| = \sqrt{0^2 + 0^2 + 0^2 + 1^2 + 1^2 + 1^2 + 1^2 + 0^2} = \sqrt{4} = 2 \\]

Substituting the computed values into the formula, we obtain a cosine similarity of:

\\[C_s(V_x, V_y) = \frac{3}{2 \times 2} = \frac{3}{4} = 0.75 \\]

## The angle between vectors

To find the angle \\(\theta\\) between the two vectors \\(V_x\\) and \\(V_y\\) from the cosine similarity value, we apply the [arccosine](<../arcsine-and-arccosine>) function. Since cosine similarity is defined as the cosine of the angle between the two vectors, the angle can be recovered by inverting that relationship. In the present case we obtain:

\\[\theta = \arccos(0.75) \approx 41.4^\circ \\]

This result is consistent with the high similarity value computed earlier: an angle of approximately \\(41.4^\circ\\) indicates that the two vectors are oriented in nearly the same direction in the feature space. In general, as the angle between two vectors decreases towards zero, their cosine similarity approaches \\(1\\), reflecting an increasing degree of similarity between the corresponding documents.

## Cosine similarity and vector orthogonality

Let us now compute the cosine similarity between \\(V_x\\) and \\(V_z\\), which represent sentences with entirely different subject matter. The dot product \\(V_x \cdot V_z\\) is obtained by multiplying the corresponding components of the two vectors and summing the results:

\\[V_x \cdot V_z = (0 \cdot 1) + (1 \cdot 0) + (0 \cdot 1) + (1 \cdot 0) + (0 \cdot 0) + (1 \cdot 0) + (1 \cdot 0) + (0 \cdot 1) = 0 \\]

Since the numerator of the cosine similarity formula is equal to zero, the cosine similarity between \\(V_x\\) and \\(V_z\\) is itself equal to zero. This result indicates that the two vectors are orthogonal, meaning they share no common features whatsoever. From a linguistic standpoint, this is entirely consistent with the observation that sentences \\(x\\) and \\(z\\) have no words in common once the stop words have been removed.

## Python implementation

Below is an example of Python code for calculating the cosine similarity of vectors \\(V_x\\) and \\(V_y\\) that you can [test](https://trinket.io/python3/cd31924b04) on an online IDE.
    
    
    import numpy as np
    # Define the vectors
    Vx = np.array([0, 1, 0, 1, 0, 1, 1, 0])
    Vy = np.array([0, 0, 0, 1, 1, 1, 1, 0])
    # Function to calculate cosine similarity
    def cosine_similarity(vector1, vector2):
        # Calculate the dot product
        dot_product = np.dot(vector1, vector2)
        # Calculate the norms of each vector
        norm1 = np.linalg.norm(vector1)
        norm2 = np.linalg.norm(vector2)
        # Calculate the cosine similarity
        cosine_sim = dot_product / (norm1 * norm2)
        return cosine_sim
    # Calculate and print the cosine similarity
    similarity = cosine_similarity(Vx, Vy)
    print(f"Cosine similarity (Vx,Vy): {similarity}")
    

> Having removed some words from the original example sentences to make the evaluation more precise, the sentences have been directly inserted as vectors in the code, returning a cosine similarity value of \\(0.75\\). For a complete example that starts from the raw sentences, the following code can be considered instead. In this case, the cosine similarity value between \\(V_x\\) and \\(V_y\\) will be approximately \\(0.48\\).

    
    
    from sklearn.feature_extraction.text import TfidfVectorizer
    from sklearn.metrics.pairwise import cosine_similarity
    sentences = [
        "I am fond of reading thriller novels.",  # x
        "I prefer reading thriller novels.",      # y
        "Yesterday, I arrived late."              # z
    ]
    # Initialize a TF-IDF Vectorizer
    vectorizer = TfidfVectorizer()
    # Fit and transform the sentences to a TF-IDF matrix
    tfidf_matrix = vectorizer.fit_transform(sentences)
    # Compute the cosine similarity matrix
    cosine_sim = cosine_similarity(tfidf_matrix, tfidf_matrix)
    print("Cosine Similarity Matrix:\n", cosine_sim)
    

## Conclusions

Cosine similarity is a mathematically well-founded and computationally efficient measure for assessing the degree of similarity between documents of moderate length. Its main strength lies in the normalization it performs: by focusing on the direction of the vectors rather than their magnitude, it correctly identifies as similar those documents that share the same thematic content, regardless of how long they are.

For longer documents, however, the situation is more nuanced. As the length of a text increases, so does the complexity of its semantic structure, and a single angle in a term-frequency vector space may no longer be sufficient to capture subtle differences in meaning. In such cases, more sophisticated approaches, such as dense vector representations produced by neural language models, are generally more appropriate.

Technique

The structure of the entry is shown in the conceptual map, where each branch represents a core component and the sub-nodes highlight the specific notions discussed.

Advanced

3

Requires

0

Enables

The following concepts, [Cosine Function](https://algebrica.org/cosine-function/), [Vector Spaces](https://algebrica.org/vector-spaces/), [Vectors](https://algebrica.org/vectors/), are required as prerequisites for this entry.
